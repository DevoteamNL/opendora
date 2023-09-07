# devlake-go

**group-sync** will retrieve all Backstage entities of type `Group` and insert the names as new teams into the `teams` table of DevLake. Duplicate team names will be skipped.

By default, the script will look for DevLake at http://localhost:4000/ and Backstage at http://localhost:7007/ but this can be changed using environment variables, respectively: `BACKSTAGE_URL` and `DEVLAKE_URL`. **Make sure to include the trailing forward slash `/`**

### TODO: 
- Parent/child group/team structure.
- Error handling in case the teams table has not-yet been created