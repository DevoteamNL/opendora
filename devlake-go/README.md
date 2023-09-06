# devlake-go

**group-sync** will retrieve all Backstage entities of type `Group` and insert the names as new teams into the `teams` table of DevLake. Duplicate team names will be skipped.

### TODO: 
- Parent/child group/team structure.
- Server URLs from environment variables
- Error handling in case the teams table has not-yet been created