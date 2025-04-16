from enum import Enum

class ValidStatuses(Enum):
    PENDING = "pending"
    COMPLETED = "completed"
    FAILED = "failed"