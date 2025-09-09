
This is an experiment to see if I was able to reproduce [Harper Reed's process](https://harper.blog/2025/05/08/basic-claude-code/ "Basic Claude Code | Harper Reed's Blog").

I created the spec and the [plan](./spec.md) with ChatGPT, then developed initially with Claude Code. When I started hitting the 5-hours limit (I'm on the $20 plan), I used other tools.  I used briefly Gemini CLI, then ONA, then Codex extensively.

Notable points

 - The acceptance tests are the strongest guardrails. They provide assurance that the code does what it's supposed to do in a clear way (check out the [acceptance tests](./acceptance/testdata)).  Unit tests are useful too, as design feedback to the AI, but they are not as straightforward to use for checking regressions and progress 
 - There is [exactly one commit](https://github.com/xpmatteo/basic-with-claude/commit/60ed4c8355bba53efb6a6ed28aff978ddcfb38ae "Human coding: show that we can track the currentLine in the ParseProg… · xpmatteo/basic-with-claude@60ed4c8 · GitHub") where I wrote code with my own hands, and it was where I was trying to get the AI to simplify the code.  The actual idea was suggested by Claude in a conversation, where he said that error line tracking is only needed during parsing, not during program execution.  I tried to get both Claude Sonnet 4 and GPT 5 to make the simple change to track the file line number in the ParseProgram method, but wasn't able to make myself understood.  In 10 minutes of human code, I made the change that proved that the line number tracking in the AST nodes wasn't needed, and then got Gemini to clean up the dead code.
 - Alternate step implementation chats with refactoring chats.  Discuss refactoring options with the AI.  Avail the awesome [Kent Beck Code Mentor](https://gist.githubusercontent.com/gscalzo/c4574328e3c211f85f8d0afc247f07c5/raw/f4b96e78f9e5581e2ce57471eb1440caa037e55d/kent-beck-code-mentor.md) agent by Giordano Scalzo.
 - Keep the project documentation up to date. Follow the [advice from Shrivu](https://blog.sshh.io/p/ai-cant-read-your-docs "Shrivus's Substack – AI Can't Read Your Docs").


Status: mostly done. I don't expect to do much work on this anymore

Check out classic basic programs such as

    scripts/run.sh testdata/guess_number.bas

or

    scripts/run.sh testdata/hamurabi.bas
    scripts/run.sh -max-steps 100000 testdata/wumpus.bas 
