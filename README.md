# Tibber DML Project

## Project Doel
Dit project is ontworpen om de streaming functionaliteit te beheren voor de `energiegemeenschap` database, specifiek voor het `tibber` schema. Het project werkt samen met het [tibber_ddl project](https://github.com/herbrandhofker/tibber_ddl) om de data streaming van Tibber API data te activeren en te beheren.

## Functionaliteit
- Het project luistert naar events die worden gegenereerd door het tibber_ddl project
- Deze events bevatten:
  - Tibber API token
  - Een `is_active` vlag
- Wanneer een event wordt ontvangen, wordt het `tibber` schema in de `energiegemeenschap` database automatisch gevuld met de relevante data

## Integratie
Dit project is een aanvulling op het [tibber_ddl project](https://github.com/herbrandhofker/tibber_ddl), dat verantwoordelijk is voor de database structuur. Terwijl het DDL project de database schema's definieert, zorgt dit DML project voor de daadwerkelijke data streaming en het vullen van de database.

## Vereisten
- Toegang tot de `energiegemeenschap` database
- Tibber API token
- Verbinding met het tibber_ddl project

## Installatie
[Installatie instructies volgen]

## Gebruik
[Gebruik instructies volgen]

## Configuratie
[Configuratie instructies volgen]

## Licentie
[Licentie informatie volgen] 