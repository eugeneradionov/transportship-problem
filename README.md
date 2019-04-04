## Heroku app: https://golang-transportship-problem.herokuapp.com/

```
coverage: 91.6% of statements
```

```
POST https://golang-transportship-problem.herokuapp.com/solve
{
    "suppliers": [
        {
            "id": 1,
            "stock": 30
        },
        {
            "id": 2,
            "stock": 40
        },
        {
            "id": 3,
            "stock": 20
        }
    ],
    "consumers": [
        {
            "id": 1,
            "demand": 20
        },
        {
            "id": 2,
            "demand": 30
        },
        {
            "id": 3,
            "demand": 30
        },
        {
            "id": 4,
            "demand": 10
        }
    ],
    "transport_cost": [
        [
            2,
            3,
            2,
            4
        ],
        [
            3,
            2,
            5,
            1
        ],
        [
            4,
            3,
            2,
            6
        ]
    ]
}
```
