Output:
```
definition user {}

definition role {
        relation can_read_things: user:*
        relation can_write_things: user:*
        relation things_admin: user:*
        relation global_admin: user:*
}

definition role_binding {
        relation subject: user
        relation role: role
        permission read = subject & role->can_read_things + role->things_admin + role->global_admin
}
```
