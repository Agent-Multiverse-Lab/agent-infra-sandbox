import json
import os
import platform


print(
    json.dumps(
        {
            "message": "hello from sandbox workload",
            "pid": os.getpid(),
            "python": platform.python_version(),
        },
        sort_keys=True,
    )
)
