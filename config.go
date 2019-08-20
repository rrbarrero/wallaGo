package main

type Configuration struct {
	DEBUG                         bool
	URL_TPLE                      string
	SCORING_VALIDATION            bool
	SCORING_VALIDATION_MIN_STARTS int8
	MIN_RECEIVED_REVIEWS_COUNT    int32
}
