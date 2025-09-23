# gradr
Lightweight CLI application to keep track of my grades for my self taught journey. Written in Golang.

Displays Class Grades

| Class                              | Grade |
|------------------------------------|-------|
| Beginning and Intermediate Algebra |  97%  |
| Geometry                           |  98%  |
| Calculus                           |  92%  |
| Linear Algebra                     |  90%  |

OR 

| Class                              | Grade |
|------------------------------------|-------|
| NO CLASSES                         |       |

Prompts "Choose an option:"
    - Lists: "Add a Class", "Select a Class", "View Overall Grades"

If they chose "Add a Class"
    - Capture Input for Class name

If they chose "Select a Class"
    - List all classes, and allow them to choose one

If they choose a class
    - Prompts "Choose an option:"
        - Lists: "Add assignment grade", "Edit assignment", "Edit assignment weights"

If they choose "Add assignment grade"
    - Capture several inputs for all the needed info

If they choose "Edit assignment"
    - Prompts "Choose an assignment to edit:"
        - Lists all assignments
        - Chooses assignment
        - Prompts "What would you like to edit:"
          - Lists everything you can edit
          - Chooses what to edit
          - Captures input for the new value

If they choose "Edit assignment weights"
    - Prompts "Choose an assignment type to edit:"
      - Lists all assignment types along with their weights
      - Choose what to edit
      - Captures input for new value

