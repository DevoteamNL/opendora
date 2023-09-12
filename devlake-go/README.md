# devlake-go

## group-sync

This will retrieve all Backstage entities of type `Group` and insert the names as new teams into the `teams` table of DevLake. This also updates the `parentId` for teams where the `childOf` and `parentOf` relationship points to an existing group/team in Backstage or DevLake.

**Note:** group names are will be saved with the correct case in DevLake but group relationships are case insensitive. If you have groups: `groupa` and `GroupA`, one of them parenting `groupb`, `groupb` will be parented to the first group found in DevLake (order not guaranteed). I recommend to just make group names unique ignoring case.

By default, the script will look for DevLake at http://localhost:4000/ and Backstage at http://localhost:7007/ but this can be changed using environment variables, respectively: `BACKSTAGE_URL` and `DEVLAKE_URL`. **Make sure to include the trailing forward slash `/`**

If you wish to replace the DevLake teams table fully, you can set the environment variable `REPLACE_DEVLAKE_TEAMS` to any value. **Be careful:** This will delete all current DevLake teams.
