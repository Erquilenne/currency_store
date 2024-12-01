package repository

const getDates = `
        SELECT DISTINCT time
        FROM currency
        ORDER BY time DESC
    `

const getByDate = `
        SELECT currency, type, value
        FROM currency
        WHERE DATE(time) = DATE($1)
    `

const save = `
        INSERT INTO currency (currency, type, value)
        VALUES ($1, $2, $3)
    `
