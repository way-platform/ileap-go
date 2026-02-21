# Add omitempty to pointer fields missing it
/\*.*`json:"[^,]*"`/s/`json:"\([^"]*\)"`/`json:"\1,omitempty"`/
