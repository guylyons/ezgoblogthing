---
title: "Drupal Links"
date: 2026-02-20
draft: false
description: "Quick links for Drupal APIs, docs, and daily CLI cheat sheets for Drush and DDEV."
---

## Core Drupal Resources

- [Drupal.org Documentation](https://www.drupal.org/docs)
- [Drupal API Reference](https://api.drupal.org/)
- [Drupal Core Release Notes](https://www.drupal.org/project/drupal/releases)
- [Drupal Security Advisories](https://www.drupal.org/security)
- [Issue Queue Search (Core)](https://www.drupal.org/project/issues/search/drupal)
- [Change Records](https://www.drupal.org/list-changes/drupal)

## Frontend and Theming

- [Twig Documentation](https://twig.symfony.com/doc/3.x/)
- [Drupal Twig in Templates](https://www.drupal.org/docs/theming-drupal/twig-in-drupal)
- [Theme System Overview](https://www.drupal.org/docs/theming-drupal)
- [Render API Overview](https://www.drupal.org/docs/drupal-apis/render-api)

## Module and Site Building APIs

- [Form API](https://api.drupal.org/api/drupal/elements)
- [Plugin API](https://www.drupal.org/docs/drupal-apis/plugin-api)
- [Entity API](https://www.drupal.org/docs/drupal-apis/entity-api)
- [Configuration API](https://www.drupal.org/docs/drupal-apis/configuration-api)
- [Routing and Controllers](https://www.drupal.org/docs/drupal-apis/routing-system)

# Commands

## Database and Tools

- `ddev import-db --src=./db.sql.gz`
- `ddev export-db --file=./db.sql.gz`
- `ddev ssh`
- `ddev logs -s web`
- `ddev xdebug on` / `ddev xdebug off`

## Composer Shortcuts for Drupal

- `composer require drupal/<module_name>`
- `composer update drupal/core-* --with-all-dependencies`
- `composer why-not drupal/core <version>`
- `composer audit`
