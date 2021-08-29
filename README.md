ClueBot NG - Review Interface
=============================

The review interface handles (trusted) human review of filtered reports & false positive/negative/random edits.

It is designed to be the main source of data for training the [https://github.com/cluebotng/core](core) ANN.

## Runtime Configuration

A dedicated local MySQL database is a hard runtime dependency.

All details are contained within `config.yaml`, which should be considered sensitive.

## Scheduled endpoints
* /api/cron/stats - Update the Wikipedia user stats page
* /api/report/import - Import report entries marked for review
* /api/report/export - Called by the report interface to update entries in review

## Training endpoints
* /api/export/done - All completed edits formatted as XML
* /api/export/dump - All edits formatted as XML
