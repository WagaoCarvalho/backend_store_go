CREATE TABLE IF NOT EXISTS clients_cpf_contact_relations (
    contact_id INT NOT NULL,
    client_cpf_id INT NOT NULL,

    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_clients_cpf_contact_relations
        PRIMARY KEY (contact_id, client_cpf_id),

    CONSTRAINT fk_clients_cpf_contact_relations_contact
        FOREIGN KEY (contact_id)
        REFERENCES contacts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_clients_cpf_contact_relations_client
        FOREIGN KEY (client_cpf_id)
        REFERENCES clients_cpf(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_clients_cpf_contact_relations_client
    ON clients_cpf_contact_relations (client_cpf_id);

CREATE INDEX idx_clients_cpf_contact_relations_contact
    ON clients_cpf_contact_relations (contact_id);
