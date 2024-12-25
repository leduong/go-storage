SHELL := /bin/bash

commit:
	git commit -m "🍻 Updated at `date`"

pull:
	git pull origin main

push:
	git push --set-upstream origin main

