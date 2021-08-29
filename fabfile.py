import time
from pathlib import PosixPath

import requests
from fabric import Connection, Config, task
from patchwork import files


def _get_latest_github_release(org, repo):
    """Return the latest release tag from GitHub"""
    r = requests.get(f"https://api.github.com/repos/{org}/{repo}/releases/latest")
    r.raise_for_status()
    return r.json()["tag_name"]


REVIEW_RELEASE = _get_latest_github_release('cluebotng', 'reviewng')
TOOL_DIR = PosixPath('/data/project/cluebotng-review')

c = Connection(
    'login.tools.wmflabs.org',
    config=Config(
        overrides={'sudo': {'user': 'tools.cluebotng-review', 'prefix': '/usr/bin/sudo -ni'}}
    ),
)


def _update():
    """Update the review release."""
    print(f'Moving reviewng to {REVIEW_RELEASE}')
    release_dir = TOOL_DIR / "releases"
    if not files.exists(c, release_dir):
        c.sudo(f'mkdir -p {release_dir}')

    target_file = release_dir / REVIEW_RELEASE
    existing_release = files.exists(c, f'{target_file.as_posix()}')
    if not existing_release:
        c.sudo(f'wget -O {target_file.as_posix()}'
               ' https://github.com/cluebotng/reviewng/releases'
               f'/download/{REVIEW_RELEASE}/reviewng')
        c.sudo(f'chmod 550 {target_file.as_posix()}')

    c.sudo(f'ln -sf {target_file.as_posix()} {TOOL_DIR / "reviewng"}')
    return not existing_release


def _update_crontab():
    mysql_dir = (TOOL_DIR / "mysql_backups")
    if not files.exists(c, mysql_dir):
        c.sudo(f'mkdir -p {mysql_dir}')

    c.sudo(f'''crontab - <<'EOL'
# Backups
45 */2 * * * /usr/bin/jsub -N cron-mysql-backup -once -quiet mysqldump --defaults-file=replica.my.cnf -h tools-db -r mysql_backups/$(date +"%d-%m-%Y_%H-%M-%S")-review.sql s54862__review
30 5 * * * /usr/bin/jsub -N cron-mysql-prune -once -quiet find mysql_backups -mtime +7 -delete

# Scheduled endpoints
13 9 * * * /usr/bin/jsub -N cron-update-stats -once -quiet curl -s https://cluebotng-review.toolforge.org/api/cron/stats
13 * * * * /usr/bin/jsub -N cron-report-import -once -quiet curl -s https://cluebotng-review.toolforge.org/api/report/import
48 * * * * /usr/bin/jsub -N cron-review-import -once -quiet curl -s https://cluebotng.toolforge.org/api/?action=review.import
EOL
    ''')


def _restart():
    c.sudo(f'webservice --backend=kubernetes golang111 restart {TOOL_DIR / "reviewng"}')


@task()
def restart(c):
    """Restart the webservice."""
    _restart()


@task()
def deploy(c):
    """Deploy the latest release & restart the webservice."""
    if _update():
        _restart()
    _update_crontab()
