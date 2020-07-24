#!/usr/bin/env python3

import sys
import re
from typing import List

from . import message

RE_MESSAGE_BREAK = re.compile(r"^-* *$")
RE_HEADER_LINE =   re.compile(r"^(?:Date|From|Subject|To|Cc|Message-ID|Content-Type):")
RE_BLANK_LINE =    re.compile(r"^ *$")
RE_SUBJECT_LINE =  re.compile(r"^Subject: *(.*) *$")
RE_DATE_LINE =     re.compile(r"^Date: *(.*) *$")
RE_FROM_LINE =     re.compile(r"^From: *(.*) *$")
RE_TO_LINE =       re.compile(r"^To: *(.*) *$")
RE_CC_LINE =       re.compile(r"^Cc: *(.*) *$")
RE_ID_LINE =       re.compile(r"^Message-ID: *(.*) *$")
RE_CONTENT_LINE =  re.compile(r"^Content-Type: *(.*) *$")
RE_RUNON =         re.compile(r"^[ \t]+(.*) *$")
RE_BOUNDARY =      re.compile(r'.*boundary="(.*)".*')

def split_messages(blob: List[bytes]) -> List[List[bytes]]:
    """Split a blob into messages."""
    message_breaks = list()
    message_start = 0
    messages = list()

    # Find probable message breaks
    for index, bytes_line in enumerate(blob):
        try:
            line = str(bytes_line)
        except:
            # test for bad encodings
            raise
        if RE_MESSAGE_BREAK.match(line):
            message_breaks.append(index)

    # Validate message breaks and copy text into split messages
    # NOTE: message breaks are validated by checking for...
    #  1) A blank line following the message break
    #  2) A header line following the blank line
    for index in message_breaks:
        try:
            line1 = str(blob[index+1])
            line2 = str(blob[index+2])
        except:
            # test for bad encodings
            raise

        # If fails validation, skip to next probable break
        if not RE_BLANK_LINE.match(line1):
            continue
        elif not RE_HEADER_LINE.match(line2):
            continue

        # Message spans from known start to line before break
        messages.append(blob[message_start:index - 1])

        # Next message starts on first header line
        message_start = index + 2

    # Handle remainder
    messages.append(blob[message_start:])

    return messages

def split_message_parts(
    blob: List[bytes],
    boundary: str,
) -> List[List[bytes]]:
    """Split a blob into message parts."""
    part_breaks = list()
    parts = list()
    part_start = 0

    # NOTE: can use `in' operator with bytes and strings
    for index, line in enumerate(blob):
        if boundary in line:
            part_breaks.append(index)

    for index in part_breaks:
        parts.append(blob[part_start:index - 1])
        part_start = index + 1

    parts.append(blob[part_start:])

    return parts





def parse_message(blob: List[bytes]):
    """Parse a message blob for metadata and parts."""
    msg = message.Message()
    header_end = 0

    # Parse the header
    for bytes_line in blob:
        header_end += 1
        try:
            line = str(bytes_line)
        except:
            # test for bad encodings
            raise

        if RE_BLANK_LINE.match(line):
            break
        elif match := RE_RUNON.match(line):
            msg.append_last(line)

        elif match := RE_SUBJECT_LINE.match(line):
            msg.hdr_subject = match.group(1)
        elif match := RE_DATE_LINE.match(line):
            msg.hdr_date = match.group(1)
        elif match := RE_FROM_LINE.match(line):
            msg.hdr_from = match.group(1)
        elif match := RE_TO_LINE.match(line):
            msg.hdr_to = match.group(1)
        elif match := RE_CC_LINE.match(line):
            msg.hdr_cc = match.group(1)
        elif match := RE_ID_LINE.match(line):
            msg.hdr_message_id = match.group(1)
        elif match := RE_CONTENT_LINE.match(line):
            msg.content_type = match.group(1)

    # store content
    msg._content = blob[header_end + 1:]

    # parse for parts
    try:
        if "multipart/" in msg.content_type:
            msg = msg.into_multipart()
    except:
        pass
    if isinstance(msg,message.MultipartMessage):
        if match := RE_BOUNDARY.match(msg.content_type):
            boundary = "".join(match.group(1).split())
        else:
            raise ValueError("no boundary for multipart content") from None
        msg._parts = split_message_parts(msg._content, boundary)

    print(msg)

def read_input() -> List[bytes]:
    """Read STDIN into a blob."""
    textblob = list()
    for line in sys.stdin:
        try:
            textblob.append(line.rstrip())
        except:
            # test for bad encoding
            raise
    return textblob

