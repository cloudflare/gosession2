package proto

// LatestOffset informs the broker that
// the consumer wishes to fetch the most
// recent event in the topic.
const LatestOffset = -1

// EarliestOffset informs the broker that
// the consumer wishes to fetch the oldest
// event in the topic.
const EarliestOffset = -2
