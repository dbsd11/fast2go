#!/bin/sh
  kill -9 $(ps -ax|grep go|grep -v grep|cut -c 1-5)
