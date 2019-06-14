package clickhouse

const ddlQuery = `
	CREATE TABLE IF NOT EXISTS events (
		EventType		LowCardinality(String) -- Enum can be used
		, Ts			DateTime
		, StringParams Nested (
			Name			LowCardinality(String)
			, Value 	String
		)
		, IntParams Nested (
			Name			LowCardinality(String)
			, Value 		Int64
		)
		, FloatParams Nested (
			Name			LowCardinality(String)
			, Value 		Float64
		)
		-- Predefined params, for advanced analytics
		, IP 			IPv4
		, UID			UUID
		-- , AppID			UUID
	) Engine MergeTree PARTITION BY toYYYYMM(Ts) ORDER BY (Ts)
	`

const insertQuery = `
	INSERT INTO events (
		EventType
		, Ts
		, StringParams.Name
		, StringParams.Value
		, IntParams.Name
		, IntParams.Value
		, FloatParams.Name
		, FloatParams.Value
		, UID
	) VALUES (?, toDateTime(?), ?, ?, ?, ?, ?, ?, ?)
	`
