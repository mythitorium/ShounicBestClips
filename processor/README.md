Myth's Epic Clip Submissions Processor
===

This is a python script that processes clip submissions (.csv) and outputs a sqlite database (.db)

It:
- Ensures submitted links are valid youtube urls
- Extracts ids from full youtube links
- Allows streamlined process to fixing horrible dogshit submissions
- Exports everything to a database that's identical in structure to the file used by the website itself.

Setup
---
This script doesn't use any external modules

* Install [Python](https://www.python.org/downloads/) (Ideally Python 3.12)
* `git clone` the main repo
* Navigate to this folder
* `py main.py` to run script
* Make sure you have a csv file of clip submissions

Notes
---
The csv layout this script expects is hardcoded. Urls are in the 4th column. Usernames are in the 2nd.

