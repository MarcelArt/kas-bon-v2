## Casbin Policy Test
### Policy
```
p, admin, sraya, EKMJ, users, read
p, admin, sraya, EKMJ, users, create
p, viewer, sraya, EKMJ, users, read
g, calvin, admin, EKMJ
g, beny, viewer, EKMJ
```

### Request
```
calvin, sraya, EKMJ, users, read
calvin, sraya, EKMJ, users, create
beny, sraya, EKMJ, users, create
```