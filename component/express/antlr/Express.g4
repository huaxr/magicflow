grammar Express;

// Tokens
NUMBER: [0-9]+;
CHAR: [A-Za-z]+;

EXP: '$'*'.'CHAR*;
Comma: ' '* ',' ' '*;
LBracket: ' '* '(' ' '*;
RBracket: ' '* ')' ' '*;

EQUALS: 'equals' LBracket EXP  Comma NUMBER RBracket;
WHITESPACE: [ \r\n\t]+ -> skip;

// Rules
start : expression EOF;

expression
    : EQUALS # EQUALS
    | EXP # EXP
    | NUMBER # NUMBER
    ;

