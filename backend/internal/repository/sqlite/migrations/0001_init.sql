-- Esquema inicial do IncidentTrack.
-- Todas as colunas de texto de enum são validadas na aplicação (camada
-- dto/validator) e reforçadas aqui por CHECK — dupla barreira.

CREATE TABLE IF NOT EXISTS users (
    id            TEXT PRIMARY KEY,
    email         TEXT NOT NULL,
    name          TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TEXT NOT NULL,
    updated_at    TEXT NOT NULL
);

-- E-mail único e case-insensitive (já gravamos em lowercase, o COLLATE
-- NOCASE é cinto e suspensório).
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email
    ON users (email COLLATE NOCASE);

CREATE TABLE IF NOT EXISTS incidents (
    id          TEXT PRIMARY KEY,
    titulo      TEXT NOT NULL,
    descricao   TEXT NOT NULL DEFAULT '',
    categoria   TEXT NOT NULL CHECK (categoria IN
                  ('phishing','malware','acesso_indevido','vazamento_dados','outro')),
    severidade  TEXT NOT NULL CHECK (severidade IN
                  ('baixa','media','alta','critica')),
    status      TEXT NOT NULL CHECK (status IN
                  ('aberto','em_analise','resolvido')),
    owner_id    TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at  TEXT NOT NULL,
    updated_at  TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_incidents_owner
    ON incidents (owner_id, created_at DESC);
