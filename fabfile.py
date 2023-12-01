import requests
import time
from fabric import Connection, Config, task
from pathlib import PosixPath


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
    c.sudo(f'mkdir -p {TOOL_DIR / "releases"}')

    target_file = TOOL_DIR / "releases" / REVIEW_RELEASE
    c.sudo(f'test -f {target_file.as_posix()} || wget -O {target_file.as_posix()}'
           f' https://github.com/cluebotng/reviewng/releases/download/{REVIEW_RELEASE}/reviewng')
    c.sudo(f'chmod 550 {target_file.as_posix()}')

    c.sudo(f'ln -sf {target_file.as_posix()} {TOOL_DIR / "reviewng"}')


def _update_crontab():
    mysql_dir = (TOOL_DIR / "mysql_backups")
    c.sudo(f'mkdir -p {mysql_dir}')

    print('Clear crontab entries')
    c.sudo('crontab -r || true')

    print('Update job entries')
    c.sudo(f'''cat > {TOOL_DIR / "jobs.yaml"} <<'EOL'
---
# Backups
- name: backup-database
  command: mysqldump --defaults-file={TOOL_DIR / "replica.my.cnf"} -h tools-db -r {mysql_dir}/$(date +"%d-%m-%Y_%H-%M-%S")-review.sql s54862__review
  image: bullseye
  filelog-stdout: logs/backup_database.stdout.log
  filelog-stderr: logs/backup_database.stderr.log
  schedule: '45 */2 * * *'
  emails: none

- name: prune-backups
  command: find {mysql_dir} -mtime +7 -delete
  image: bullseye
  filelog-stdout: logs/prune_backups.stdout.log
  filelog-stderr: logs/prune_backups.stderr.log
  schedule: '30 5 * * *'
  emails: none

# Scheduled endpoints
- name: update-stats
  command: curl -s https://cluebotng-review.toolforge.org/api/cron/stats
  image: bullseye
  filelog-stdout: logs/update_stats.stdout.log
  filelog-stderr: logs/update_stats.stderr.log
  schedule: '13 9 * * *'
  emails: none

- name: report-import
  command: curl -s https://cluebotng-review.toolforge.org/api/report/import
  image: bullseye
  filelog-stdout: logs/report_import.stdout.log
  filelog-stderr: logs/report_import.stderr.log
  schedule: '13 * * * *'
  emails: none

- name: review-import
  command: curl -s https://cluebotng.toolforge.org/api/?action=review.import
  image: bullseye
  filelog-stdout: logs/review_import.stdout.log
  filelog-stderr: logs/review_import.stderr.log
  schedule: '48 * * * *'
  emails: none

- name: training-import
  command: curl -s https://cluebotng-review.toolforge.org/api/training/import
  image: bullseye
  filelog-stdout: logs/training_import.stdout.log
  filelog-stderr: logs/training_import.stderr.log
  schedule: '30 * * * *'
  emails: none
EOL
    ''')
    c.sudo(f'XDG_CONFIG_HOME={TOOL_DIR} toolforge jobs load {TOOL_DIR / "jobs.yaml"}')


def _restart():
    c.sudo(f'webservice --backend=kubernetes golang1.11 stop {TOOL_DIR / "reviewng"}')
    c.sudo(f'webservice --backend=kubernetes golang1.11 start {TOOL_DIR / "reviewng"}')


@task()
def restart(c):
    """Restart the webservice."""
    _restart()


@task()
def deploy(c):
    """Deploy the latest release & restart the webservice."""
    _update()
    _restart()
    _update_crontab()
