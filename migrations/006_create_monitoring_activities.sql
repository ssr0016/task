CREATE TABLE activity_logs(
    id SERIAL NOT NULL,
    user_id varchar(255) NOT NULL,
    activity varchar(255) NOT NULL,
    "action" varchar(255) NOT NULL,
    resource varchar(255) NOT NULL,
    details varchar(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);