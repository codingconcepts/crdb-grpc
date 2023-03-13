CREATE TABLE todo (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "title" STRING NOT NULL
);

INSERT INTO todo (id, title) VALUES
    ('8da8291c-5985-41c9-8069-0de865dd20d7', 'todo a'),
    ('03036918-ba3d-4264-84ed-94e0c6d7433e', 'todo b'),
    ('81a39960-366f-43dd-ad10-25b33454168f', 'todo c');