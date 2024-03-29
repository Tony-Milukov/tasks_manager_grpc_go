CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email VARCHAR(255) NOT NULL UNIQUE,
        password TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS statuses (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        statusId INT,
        creatorId TEXT,
        due TIMESTAMP,
        completed bool,
        FOREIGN KEY(statusId) REFERENCES statuses(id),
        FOREIGN KEY(creatorId) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS task_assignees (
        id SERIAL PRIMARY KEY,
        role VARCHAR(255) NOT NULL,
        userId TEXT,
        taskId int,
        FOREIGN KEY(userId) REFERENCES users(id),
        FOREIGN KEY(taskId) REFERENCES tasks(id),
        UNIQUE (userId, taskId)
);
