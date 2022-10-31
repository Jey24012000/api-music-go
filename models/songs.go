package models

import (
	"encoding/xml"

	"gorm.io/gorm"
)

type Response struct {
	Songs []Songs `json:"results"`
}

type Songs struct {
	gorm.Model
	ArtistId        int     `json:"artistId" gorm:"primary_key"`
	TrackName       string  `json:"trackName"`
	ArtistName      string  `json:"artistName"`
	TrackTimeMillis int     `json:"trackTimeMillis"`
	CollectionName  string  `json:"collectionName"`
	ArtworkUrl30    string  `json:"artworkUrl30"`
	TrackPrice      float32 `json:"trackPrice"`
	Origin          string
}



type GetLyricResult struct {
	XMLName xml.Name `xml:"GetLyricResult"`
	//Text              string   `xml:",chardata"`
	//Xsd               string   `xml:"xsd,attr"`
	//Xsi               string   `xml:"xsi,attr"`
	//Xmlns             string   `xml:"xmlns,attr"`
	TrackChecksum     string `xml:"TrackChecksum"`
	TrackId           string `xml:"TrackId"`
	LyricChecksum     string `xml:"LyricChecksum"`
	LyricId           int `xml:"LyricId"`
	LyricSong         string `xml:"LyricSong"`
	LyricArtist       string `xml:"LyricArtist"`
	LyricUrl          string `xml:"LyricUrl"`
	LyricCovertArtUrl string `xml:"LyricCovertArtUrl"`
	LyricRank         int `xml:"LyricRank"`
	LyricCorrectUrl   string `xml:"LyricCorrectUrl"`
	Lyric             string `xml:"Lyric"`
}
