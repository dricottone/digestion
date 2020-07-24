#!/usr/bin/env python3

"""The message object and API."""

from typing import Optional

class Message(object):
    """Container for a message and metadata."""
    def __init__(
        self,
        *,
        hdr_subject: Optional[str] = None,
        hdr_date: Optional[str] = None,
        hdr_from: Optional[str] = None,
        hdr_to: Optional[str] = None,
        hdr_cc: Optional[str] = None,
        hdr_message_id: Optional[str] = None,
        content_type: Optional[str] = None,
        content: Optional[str] = None,
    ) -> None:
        self._subject = hdr_subject
        self._date = hdr_date
        self._from = hdr_from
        self._to = hdr_to
        self._cc = hdr_cc
        self._message_id = hdr_message_id
        self._content_type = content_type
        self._content = content
        self._last_hdr = None

    def __str__(self) -> str:
        return (
            f"Subject: {self._subject}\n"
            f"Date: {self._date}\n"
            f"To: {self._to}\n"
            f"From: {self._from}\n"
            f"Cc: {self._cc}\n"
            f"Message-ID: {self._message_id}\n"
            f"Content-Type: {self._content_type}\n"
        )

    @property
    def hdr_subject(self) -> str:
        if self._subject is not None:
            return self._subject
        else:
            raise ValueError("no header `subject' set") from None
    @hdr_subject.setter
    def hdr_subject(self, value: str):
        if self._subject is None:
            self._subject = value
            self._last_hdr = "_subject"
        else:
            raise ValueError("header `subject' already set") from None

    @property
    def hdr_date(self) -> str:
        if self._date is not None:
            return self._date
        else:
            raise ValueError("no header `date' set") from None
    @hdr_date.setter
    def hdr_date(self, value: str):
        if self._date is None:
            self._date = value
            self._last_hdr = "_date"
        else:
            raise ValueError("header `date' already set") from None

    @property
    def hdr_from(self) -> str:
        if self._from is not None:
            return self._from
        else:
            raise ValueError("no header `from' set") from None
    @hdr_from.setter
    def hdr_from(self, value: str):
        if self._from is None:
            self._from = value
            self._last_hdr = "_from"
        else:
            raise ValueError("header `from' already set") from None

    @property
    def hdr_to(self) -> str:
        if self._to is not None:
            return self._to
        else:
            raise ValueError("no header `to' set") from None
    @hdr_to.setter
    def hdr_to(self, value: str):
        if self._to is None:
            self._to = value
            self._last_hdr = "_to"
        else:
            raise ValueError("header `to' already set") from None

    @property
    def hdr_cc(self) -> str:
        if self._cc is not None:
            return self._cc
        else:
            raise ValueError("no header `cc' set") from None
    @hdr_cc.setter
    def hdr_cc(self, value: str):
        if self._cc is None:
            self._cc = value
            self._last_hdr = "_cc"
        else:
            raise ValueError("header `cc' already set") from None

    @property
    def hdr_message_id(self) -> str:
        if self._message_id is not None:
            return self._message_id
        else:
            raise ValueError("no header `message_id' set") from None
    @hdr_message_id.setter
    def hdr_message_id(self, value: str):
        if self._message_id is None:
            self._message_id = value
            self._last_hdr = "_message_id"
        else:
            raise ValueError("header `message_id' already set") from None

    @property
    def content_type(self) -> str:
        if self._content_type is not None:
            return self._content_type
        else:
            raise ValueError("no `content_type' set") from None
    @content_type.setter
    def content_type(self, value: str):
        if self._content_type is None:
            self._content_type = value
            self._last_hdr = "_content_type"
        else:
            raise ValueError("`content_type' already set") from None

    def append_last(self, value: str):
        if self._last_hdr is not None:
            old = getattr(self, self._last_hdr)
            try:
                new = old + value
            except:
                # test for bad encoding
                raise
            setattr(self, self._last_hdr, new)
        else:
            raise ValueError("no header set") from None

    def into_multipart(self):
        return MultipartMessage(
            hdr_subject=self._subject,
            hdr_date=self._date,
            hdr_from=self._from,
            hdr_to=self._to,
            hdr_cc=self._cc,
            hdr_message_id=self._message_id,
            content_type=self._content_type,
            content=self._content,
        )

class MultipartMessage(Message):
    """Container for a multi-part message and metadata."""
    def __init__(
        self,
        *,
        hdr_subject: Optional[str] = None,
        hdr_date: Optional[str] = None,
        hdr_from: Optional[str] = None,
        hdr_to: Optional[str] = None,
        hdr_cc: Optional[str] = None,
        hdr_message_id: Optional[str] = None,
        content_type: Optional[str] = None,
        content: Optional[str] = None,
    ) -> None:
        self._subject = hdr_subject
        self._date = hdr_date
        self._from = hdr_from
        self._to = hdr_to
        self._cc = hdr_cc
        self._message_id = hdr_message_id
        self._content_type = content_type
        self._content = content
        self._last_hdr = None

        self._parts = list()

    def __str__(self) -> str:
        return (
            f"Subject: {self._subject}\n"
            f"Date: {self._date}\n"
            f"To: {self._to}\n"
            f"From: {self._from}\n"
            f"Cc: {self._cc}\n"
            f"Message-ID: {self._message_id}\n"
            f"Content-Type: {self._content_type}\n"
            f"Parts: {len(self._parts)}\n"
        )

