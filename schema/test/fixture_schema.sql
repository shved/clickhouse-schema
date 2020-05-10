CREATE DATABASE testdb ENGINE = Ordinary

CREATE TABLE testdb.dates (`date` Date) ENGINE = MergeTree(date, date, 8192)

CREATE TABLE testdb.payments (`id` UUID, `status` Int32, `amount_cents` Int32, `amount_currency` String, `created_at` DateTime, `type` String, `initiator` Int32, `invoice_billing_account_id` UUID, `user_username` String, `product_id` UUID, `product_name` String, `plan_id` UUID, `plan_name` String) ENGINE = MergeTree() PARTITION BY toYYYYMM(created_at) ORDER BY created_at SETTINGS index_granularity = 8192

CREATE TABLE testdb.ping_events_data_collapsed (`action` String, `app_name` String, `customer_id` String, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime, `event_date` Date, `Sign` Int8) ENGINE = CollapsingMergeTree(Sign) PARTITION BY toYYYYMM(event_date) ORDER BY (session_id, event_date, event_time, stream_id) SETTINGS index_granularity = 8192

CREATE TABLE testdb.ping_events_data_collapsed_10s (`action` String, `app_name` String, `customer_id` String, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime, `event_date` Date, `Sign` Int8) ENGINE = CollapsingMergeTree(Sign) PARTITION BY toYYYYMM(event_date) ORDER BY (session_id, event_date, event_time, stream_id) SETTINGS index_granularity = 8192

CREATE TABLE testdb.ping_events_data_raw (`action` String, `app_name` String, `customer_id` String, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime, `event_date` Date) ENGINE = MergeTree() PARTITION BY toYYYYMM(event_date) ORDER BY (session_id, event_date, event_time, stream_id) SETTINGS index_granularity = 8192

CREATE MATERIALIZED VIEW testdb.ping_events_handler_collapsed TO testdb.ping_events_data_collapsed (`action` String, `app_name` String, `customer_id` UUID, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime('Etc/UTC'), `event_date` Date, `Sign` Int8) AS SELECT action, app_name, customer_id, stream_id, session_id, version, device_id, ip_addr, resource_type, toStartOfMinute(toDateTime(event_time)) AS event_time, toDate(event_time) AS event_date, toInt8(1) AS Sign FROM testdb.ping_events_stream GROUP BY action, app_name, session_id, customer_id, version, stream_id, resource_type, event_date, event_time, device_id, ip_addr

CREATE MATERIALIZED VIEW testdb.ping_events_handler_collapsed_10s TO testdb.ping_events_data_collapsed_10s (`action` String, `app_name` String, `customer_id` UUID, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime('Etc/UTC'), `event_date` Date, `Sign` Int8) AS SELECT action, app_name, customer_id, stream_id, session_id, version, device_id, ip_addr, resource_type, toStartOfInterval(toDateTime(event_time), toIntervalSecond(10)) AS event_time, toDate(event_time) AS event_date, toInt8(1) AS Sign FROM testdb.ping_events_stream GROUP BY action, app_name, session_id, customer_id, version, stream_id, resource_type, event_date, event_time, device_id, ip_addr

CREATE MATERIALIZED VIEW testdb.ping_events_handler_raw TO testdb.ping_events_data_raw (`action` String, `app_name` String, `customer_id` UUID, `stream_id` UUID, `session_id` UUID, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` DateTime, `event_date` Date) AS SELECT action, app_name, customer_id, stream_id, session_id, version, device_id, ip_addr, resource_type, toDateTime(event_time) AS event_time, toDate(event_time) AS event_date FROM testdb.ping_events_stream

CREATE TABLE testdb.ping_events_stream (`action` String, `app_name` String, `customer_id` String, `stream_id` UUID, `session_id` UUID, `duration` Int32, `timestamp` Int32, `version` String, `device_id` String, `ip_addr` String, `resource_type` String, `event_time` UInt32, `event_date` Date) ENGINE = Null()

CREATE TABLE testdb.schema_migrations (`version` Int64, `dirty` UInt8, `sequence` UInt64) ENGINE = TinyLog

CREATE TABLE testdb.sessions_metadata (`session_uid` String, `requested_at` Int32, `user_timezone_offset` Int32, `user_uid` String, `user_registration_state` String, `user_msisdn` String, `user_gender` String, `account_uid` String, `device_uid` String, `device_ipv4` String, `device_ipv6` String, `device_connection_type` String, `device_type` String, `device_model` String, `device_is` String, `device_os_version` String, `device_vendor` String, `device_ifa` String, `country_iso_code` String, `region_iso_code` String, `maxmind2_country_uid` String, `maxmind2_region_uid` String, `maxmind2_city_uid` String, `ipgeobase_city_uid` String, `resource_uid` String, `resource_type` String, `client_version` String, `client_environment` String, `server_version` String, `server_environment` String, `user_subscription_state` String, `ads_enabled` UInt8, `ads_cooldown` UInt8, `streaming_type` String, `_sign` Int8, `right_holder_name` String, `external_catalog_name` String, `channel_name` String, `movie_name` String, `episode_name` String, `series_name` String, `season_number` String, `episode_number` String, `program_event_name` String, `program_event_start_at` DateTime, `program_event_end_at` DateTime) ENGINE = CollapsingMergeTree(_sign) PARTITION BY toYYYYMM(toDate(requested_at)) ORDER BY (session_uid, user_uid) SETTINGS index_granularity = 8192

CREATE TABLE testdb.streamer_logs_raw (`session_id` String, `body_size` Int64, `ad_p_hash` Int64, `ad_dtmf_code` String, `adm_spot_uid` String, `requested_host` String, `streamer_host` String, `original_chunk_name` String, `event_time` DateTime, `event_date` Date) ENGINE = MergeTree() PARTITION BY toYYYYMM(event_date) ORDER BY (session_id, event_date, event_time) SETTINGS index_granularity = 8192

CREATE MATERIALIZED VIEW testdb.streamer_logs_raw_handler TO testdb.streamer_logs_raw (`session_id` String, `timestamp` Int64, `body_size` Int64, `ad_p_hash` Int64, `ad_dtmf_code` String, `adm_spot_uid` String, `requested_host` String, `streamer_host` String, `original_chunk_name` String, `event_time` DateTime, `event_date` Date) AS SELECT * FROM testdb.streamer_logs_stream

CREATE TABLE testdb.streamer_logs_stream (`session_id` String, `timestamp` Int64, `body_size` Int64, `ad_p_hash` Int64, `ad_dtmf_code` String, `adm_spot_uid` String, `requested_host` String, `streamer_host` String, `original_chunk_name` String, `event_time` DateTime, `event_date` Date) ENGINE = Null()

CREATE TABLE testdb.subscriptions (`id` UUID, `day` Date, `subscription_id` UUID, `user_id` UUID, `user_username` String, `user_tags` String, `plan_id` UUID, `plan_name` String, `plan_price_list_names` String, `plan_price_list_ids` String, `product_id` UUID, `product_name` String, `phase_id` UUID, `phase_type` String, `phase_duration_value` Int32, `phase_duration_unit` String, `phase_price_cents` Int32, `phase_price_currency` String, `has_active_subscription` String, `has_active_autorenew` String, `has_access_granted` String, `created_at` DateTime, `updated_at` DateTime, `last_used_payment_method_type` String, `last_user_activity_at` DateTime, `user_external_uid` String) ENGINE = MergeTree() PARTITION BY toYYYYMM(created_at) ORDER BY (day, created_at) SETTINGS index_granularity = 8192

CREATE TABLE testdb.watch_progress (`customer_id` String, `resource_id` String, `resource_type` String, `timestamp` Int32, `duration` Int32, `event_time` DateTime, `session_id` String, `action` String) ENGINE = ReplacingMergeTree(event_time) PARTITION BY toYYYYMM(event_time) ORDER BY (customer_id, resource_id) SETTINGS index_granularity = 8192

CREATE MATERIALIZED VIEW testdb.watch_progress_handler TO testdb.watch_progress (`customer_id` String, `resource_id` String, `resource_type` String, `timestamp` Int32, `duration` Int32, `event_time` DateTime, `session_id` DateTime, `action` String) AS SELECT toString(customer_id) AS customer_id, toString(stream_id) AS resource_id, resource_type, timestamp, duration, toDateTime(event_time) AS event_time, toString(session_id) AS session_id, action FROM testdb.ping_events_stream
