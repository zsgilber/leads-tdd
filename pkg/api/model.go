package api

import (
	"database/sql"

	"github.com/lib/pq"
)

type lead struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type leadsToDelete struct {
	Ids []int `json:"idsToDelete"`
}

func (l *lead) getLead(db *sql.DB) error {
	return db.QueryRow("SELECT first_name, last_name FROM leads.leads WHERE id=$1",
		l.ID).Scan(&l.Firstname, &l.Lastname)
}

func (l *lead) updateLead(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE leads.leads SET first_name=$1, last_name=$2 WHERE id=$3",
			l.Firstname, l.Lastname, l.ID)

	return err
}

func (l *lead) deleteLead(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM leads.leads WHERE id=$1", l.ID)

	return err
}

func (l2d *leadsToDelete) deleteLeads(db *sql.DB) error {
	//_, err := db.Exec("DELETE FROM leads.leads WHERE id IN $1", leadsToDelete)
	stmt, err := db.Prepare("DELETE FROM leads.leads WHERE id=ANY($1)")
	_, err = stmt.Exec(pq.Array(l2d.Ids))

	return err
}

func (l *lead) createLead(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO leads.leads(first_name, last_name) VALUES($1, $2) RETURNING id",
		l.Firstname, l.Lastname).Scan(&l.ID)

	if err != nil {
		return err
	}

	return nil
}

func getLeads(db *sql.DB, start, count int) ([]lead, error) {
	rows, err := db.Query(
		"SELECT id, first_name, last_name FROM leads.leads LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	leads := []lead{}

	for rows.Next() {
		var l lead
		if err := rows.Scan(&l.ID, &l.Firstname, &l.Lastname); err != nil {
			return nil, err
		}
		leads = append(leads, l)
	}

	return leads, nil
}
