CREATE INDEX idx_users_team_id ON users(team_id);

CREATE INDEX idx_pull_requests_author_id ON pull_requests(author_id);

CREATE INDEX idx_assignments_user_id ON assignments(user_id);

CREATE INDEX idx_assignments_pr_id ON assignments(pr_id);
