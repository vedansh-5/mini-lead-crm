#!/bin/bash

echo "Seeding Mini Lead CRM database..."

curl -s -X POST http://localhost:8080/leads/bulk \
-H "Content-Type: application/json" \
-d '[
  {
    "name": "Alice Smith",
    "email": "alice@example.com",
    "phone": "555-0101",
    "source": "website"
  },
  {
    "name": "Bob Jones",
    "email": "bob@example.com",
    "phone": "555-0102",
    "source": "referral"
  },
  {
    "name": "Charlie Brown",
    "email": "charlie@example.com",
    "phone": "555-0103",
    "source": "campaign"
  },
  {
    "name": "Diana Prince",
    "email": "diana@example.com",
    "phone": "555-0104",
    "source": "website"
  },
  {
    "name": "Evan Wright",
    "email": "evan@example.com",
    "phone": "555-0105",
    "source": "referral"
  }
]' | jq . || echo "Seed completed (install jq for pretty output)"

echo -e "\nData seeded successfully!"
