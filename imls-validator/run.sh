#!/usr/bin/env bash
python /code/rabbit/manage.py migrate && python /code/rabbit/manage.py runserver 0.0.0.0:8000
