from dataclasses import dataclass
from typing import List


@dataclass
class TorrentInfo:
    title: str
    href: str
    forum: str
    author: str
    size: str
    seeds: str
    leeches: str
    downloads:str
    added_date: str

@dataclass
class TorrentInfoResponse:
    page: int
    total_pages: int
    total_results: int
    results: List[TorrentInfo]

@dataclass
class MagnetInfo:
  magnet: str
  hash: str