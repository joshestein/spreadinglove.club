CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pending_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    status TEXT DEFAULT 'pending' CHECK(status IN ('pending', 'approved', 'rejected')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO messages (id, content) VALUES
    (1, 'Thank you for being.'),
    (2, 'You are worthy just as you are.'),
    (3, 'You are enough.'),
    (4, 'Do not believe the lies you tell yourself.'),
    (5, 'You are loved.'),
    (6, 'I am proud of you.'),
    (7, "If you really want something, you will find a way. If you don't, you will find an excuse."),
    (8, 'What is stopping you?'),
    (9, 'Love with all your heart.'),
    (10, 'See the world with eyes of love, not fear.'),
    (11, 'You are allowed to rest.'),
    (12, 'I am grateful you are here.'),
    (13, 'Treat yourself like someone you love.');
