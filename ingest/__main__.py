#!/usr/bin/env python3

from . import parse

for msg in parse.split_messages(parse.read_input())[1:]:
    try:
        parse.parse_message(msg)
    except:
        raise


