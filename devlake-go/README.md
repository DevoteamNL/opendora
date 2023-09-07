# devlake-go

## group-sync
This will retrieve all Backstage entities of type `Group` and insert the names as new teams into the `teams` table of DevLake. This also updates the `parentId` for teams where the `childOf` and `parentOf` relationship points to an existing group/team in Backstage or DevLake.

**Note:** group names are will be saved with the correct case in DevLake but group relationships are case insensitive. If you have groups: `groupa` and `GroupA`, one of them parenting `groupb`, `groupb` will be parented to the first group found in DevLake (order not guaranteed). I recommend to just make group names unique ignoring case.

### TODO:
- Server URLs from environment variables
- Error handling in case the teams table has not-yet been created