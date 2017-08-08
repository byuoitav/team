#!/usr/bin/env python

import os
import stat

byuoitavdir = os.environ['GOPATH']+"/src/github.com/byuoitav/"
#find all of the repositories in byuoitav/
repos = os.listdir(byuoitavdir)

for a in repos:
    if not os.path.isdir(byuoitavdir + "/" + a):
        continue
    if not os.path.exists(byuoitavdir + "/" + a + "/.git"):
        continue
    print("doing " + str(a))
    #delete the old link if it exists
    if os.path.isfile(byuoitavdir + "/" + a + "/.git/hooks/pre-commit"):
        print("removing old symlink")
        os.remove(byuoitavdir + "/" + a + "/.git/hooks/pre-commit")
    #create the link
    print("creating symlink")
    os.symlink(byuoitavdir + "/team/tools/version_hook", byuoitavdir + "/" + a + "/.git/hooks/pre-commit")
    print("setting execution rights")
    st = os.stat(byuoitavdir + "/" + a + "/.git/hooks/pre-commit")
    os.chmod( byuoitavdir + "/" + a + "/.git/hooks/pre-commit", st.st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)

