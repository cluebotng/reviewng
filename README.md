ClueBot NG - Review Interface
=============================

The review interface handles (trusted) human review of filtered reports & false positive/negative/random edits.

It is designed to be the main source of data for training the [https://github.com/cluebotng/core](core) ANN.

## Runtime Configuration

A dedicated local MySQL database is a hard runtime dependency.

All details are contained within `config.yaml`, which should be considered sensitive.
