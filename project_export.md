# Export de projet

_G√©n√©r√© le 2025-12-25T19:28:03+01:00_

## .env

```text
# Infrastructure
REDIS_URL=redis://redis:6379

# RabbitMQ - URLs pour chaque service (Vhost: benchmarks)
AMQP_URL_LEADER=amqp://bench_leader:qsd65f4c98dc7fd9s87ga6fsd5g4zsdrf9g879dfs7g@192.168.1.12:5672/benchmarks
AMQP_URL_GO=amqp://bench_go:9q8s7d9qs87dqs654dq6s54d6qs54dqs321dqs2d1qs98d7qs9d8q7@192.168.1.12:5672/benchmarks
AMQP_URL_RUST=amqp://bench_rust:qs9d87f+qsdf87qsd98f7sd654f1qsd32f1sd56f4sd@192.168.1.12:5672/benchmarks
AMQP_URL_PYTHON=amqp://bench_python:c1s9d8sf4d9s8d7fs7fqs98d6f8qs7d6f98qsd@192.168.1.12:5672/benchmarks
AMQP_URL_NODE=amqp://bench_node:sq9d8f7sd98f7sd65f4sd6f54sdf654@192.168.1.12:5672/benchmarks
```

## .git/COMMIT_EDITMSG

```text
worker-go: publier metrics et enrichir metadata dans handler fibonacci (√©viter QueueDeclare redondant)

```

## .git/FETCH_HEAD

```text
9dcf46457a2239e54fbe4a9f00f2826996522e04		branch 'master' of github.com:sicDANGBE/fibo

```

## .git/HEAD

```text
ref: refs/heads/master

```

## .git/ORIG_HEAD

```text
c9a83108423860e7c4e872bdb748dfb3d2f60be9

```

## .git/config

```text
[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
[remote "origin"]
	url = git@github.com:sicDANGBE/fibo.git
	fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
	remote = origin
	merge = refs/heads/master
	gk-last-accessed = 2025-12-25T17:17:54.154Z
	gk-last-modified = 2025-12-25T17:17:54.154Z

```

## .git/description

```text
Unnamed repository; edit this file 'description' to name the repository.

```

## .git/hooks/applypatch-msg.sample

```text
#!/bin/sh
#
# An example hook script to check the commit log message taken by
# applypatch from an e-mail message.
#
# The hook should exit with non-zero status after issuing an
# appropriate message if it wants to stop the commit.  The hook is
# allowed to edit the commit message file.
#
# To enable this hook, rename this file to "applypatch-msg".

. git-sh-setup
commitmsg="$(git rev-parse --git-path hooks/commit-msg)"
test -x "$commitmsg" && exec "$commitmsg" ${1+"$@"}
:

```

## .git/hooks/commit-msg.sample

```text
#!/bin/sh
#
# An example hook script to check the commit log message.
# Called by "git commit" with one argument, the name of the file
# that has the commit message.  The hook should exit with non-zero
# status after issuing an appropriate message if it wants to stop the
# commit.  The hook is allowed to edit the commit message file.
#
# To enable this hook, rename this file to "commit-msg".

# Uncomment the below to add a Signed-off-by line to the message.
# Doing this in a hook is a bad idea in general, but the prepare-commit-msg
# hook is more suited to it.
#
# SOB=$(git var GIT_AUTHOR_IDENT | sed -n 's/^\(.*>\).*$/Signed-off-by: \1/p')
# grep -qs "^$SOB" "$1" || echo "$SOB" >> "$1"

# This example catches duplicate Signed-off-by lines.

test "" = "$(grep '^Signed-off-by: ' "$1" |
	 sort | uniq -c | sed -e '/^[ 	]*1[ 	]/d')" || {
	echo >&2 Duplicate Signed-off-by lines.
	exit 1
}

```

## .git/hooks/fsmonitor-watchman.sample

```text
#!/usr/bin/perl

use strict;
use warnings;
use IPC::Open2;

# An example hook script to integrate Watchman
# (https://facebook.github.io/watchman/) with git to speed up detecting
# new and modified files.
#
# The hook is passed a version (currently 2) and last update token
# formatted as a string and outputs to stdout a new update token and
# all files that have been modified since the update token. Paths must
# be relative to the root of the working tree and separated by a single NUL.
#
# To enable this hook, rename this file to "query-watchman" and set
# 'git config core.fsmonitor .git/hooks/query-watchman'
#
my ($version, $last_update_token) = @ARGV;

# Uncomment for debugging
# print STDERR "$0 $version $last_update_token\n";

# Check the hook interface version
if ($version ne 2) {
	die "Unsupported query-fsmonitor hook version '$version'.\n" .
	    "Falling back to scanning...\n";
}

my $git_work_tree = get_working_dir();

my $retry = 1;

my $json_pkg;
eval {
	require JSON::XS;
	$json_pkg = "JSON::XS";
	1;
} or do {
	require JSON::PP;
	$json_pkg = "JSON::PP";
};

launch_watchman();

sub launch_watchman {
	my $o = watchman_query();
	if (is_work_tree_watched($o)) {
		output_result($o->{clock}, @{$o->{files}});
	}
}

sub output_result {
	my ($clockid, @files) = @_;

	# Uncomment for debugging watchman output
	# open (my $fh, ">", ".git/watchman-output.out");
	# binmode $fh, ":utf8";
	# print $fh "$clockid\n@files\n";
	# close $fh;

	binmode STDOUT, ":utf8";
	print $clockid;
	print "\0";
	local $, = "\0";
	print @files;
}

sub watchman_clock {
	my $response = qx/watchman clock "$git_work_tree"/;
	die "Failed to get clock id on '$git_work_tree'.\n" .
		"Falling back to scanning...\n" if $? != 0;

	return $json_pkg->new->utf8->decode($response);
}

sub watchman_query {
	my $pid = open2(\*CHLD_OUT, \*CHLD_IN, 'watchman -j --no-pretty')
	or die "open2() failed: $!\n" .
	"Falling back to scanning...\n";

	# In the query expression below we're asking for names of files that
	# changed since $last_update_token but not from the .git folder.
	#
	# To accomplish this, we're using the "since" generator to use the
	# recency index to select candidate nodes and "fields" to limit the
	# output to file names only. Then we're using the "expression" term to
	# further constrain the results.
	my $last_update_line = "";
	if (substr($last_update_token, 0, 1) eq "c") {
		$last_update_token = "\"$last_update_token\"";
		$last_update_line = qq[\n"since": $last_update_token,];
	}
	my $query = <<"	END";
		["query", "$git_work_tree", {$last_update_line
			"fields": ["name"],
			"expression": ["not", ["dirname", ".git"]]
		}]
	END

	# Uncomment for debugging the watchman query
	# open (my $fh, ">", ".git/watchman-query.json");
	# print $fh $query;
	# close $fh;

	print CHLD_IN $query;
	close CHLD_IN;
	my $response = do {local $/; <CHLD_OUT>};

	# Uncomment for debugging the watch response
	# open ($fh, ">", ".git/watchman-response.json");
	# print $fh $response;
	# close $fh;

	die "Watchman: command returned no output.\n" .
	"Falling back to scanning...\n" if $response eq "";
	die "Watchman: command returned invalid output: $response\n" .
	"Falling back to scanning...\n" unless $response =~ /^\{/;

	return $json_pkg->new->utf8->decode($response);
}

sub is_work_tree_watched {
	my ($output) = @_;
	my $error = $output->{error};
	if ($retry > 0 and $error and $error =~ m/unable to resolve root .* directory (.*) is not watched/) {
		$retry--;
		my $response = qx/watchman watch "$git_work_tree"/;
		die "Failed to make watchman watch '$git_work_tree'.\n" .
		    "Falling back to scanning...\n" if $? != 0;
		$output = $json_pkg->new->utf8->decode($response);
		$error = $output->{error};
		die "Watchman: $error.\n" .
		"Falling back to scanning...\n" if $error;

		# Uncomment for debugging watchman output
		# open (my $fh, ">", ".git/watchman-output.out");
		# close $fh;

		# Watchman will always return all files on the first query so
		# return the fast "everything is dirty" flag to git and do the
		# Watchman query just to get it over with now so we won't pay
		# the cost in git to look up each individual file.
		my $o = watchman_clock();
		$error = $output->{error};

		die "Watchman: $error.\n" .
		"Falling back to scanning...\n" if $error;

		output_result($o->{clock}, ("/"));
		$last_update_token = $o->{clock};

		eval { launch_watchman() };
		return 0;
	}

	die "Watchman: $error.\n" .
	"Falling back to scanning...\n" if $error;

	return 1;
}

sub get_working_dir {
	my $working_dir;
	if ($^O =~ 'msys' || $^O =~ 'cygwin') {
		$working_dir = Win32::GetCwd();
		$working_dir =~ tr/\\/\//;
	} else {
		require Cwd;
		$working_dir = Cwd::cwd();
	}

	return $working_dir;
}

```

## .git/hooks/post-update.sample

```text
#!/bin/sh
#
# An example hook script to prepare a packed repository for use over
# dumb transports.
#
# To enable this hook, rename this file to "post-update".

exec git update-server-info

```

## .git/hooks/pre-applypatch.sample

```text
#!/bin/sh
#
# An example hook script to verify what is about to be committed
# by applypatch from an e-mail message.
#
# The hook should exit with non-zero status after issuing an
# appropriate message if it wants to stop the commit.
#
# To enable this hook, rename this file to "pre-applypatch".

. git-sh-setup
precommit="$(git rev-parse --git-path hooks/pre-commit)"
test -x "$precommit" && exec "$precommit" ${1+"$@"}
:

```

## .git/hooks/pre-commit.sample

```text
#!/bin/sh
#
# An example hook script to verify what is about to be committed.
# Called by "git commit" with no arguments.  The hook should
# exit with non-zero status after issuing an appropriate message if
# it wants to stop the commit.
#
# To enable this hook, rename this file to "pre-commit".

if git rev-parse --verify HEAD >/dev/null 2>&1
then
	against=HEAD
else
	# Initial commit: diff against an empty tree object
	against=$(git hash-object -t tree /dev/null)
fi

# If you want to allow non-ASCII filenames set this variable to true.
allownonascii=$(git config --type=bool hooks.allownonascii)

# Redirect output to stderr.
exec 1>&2

# Cross platform projects tend to avoid non-ASCII filenames; prevent
# them from being added to the repository. We exploit the fact that the
# printable range starts at the space character and ends with tilde.
if [ "$allownonascii" != "true" ] &&
	# Note that the use of brackets around a tr range is ok here, (it's
	# even required, for portability to Solaris 10's /usr/bin/tr), since
	# the square bracket bytes happen to fall in the designated range.
	test $(git diff --cached --name-only --diff-filter=A -z $against |
	  LC_ALL=C tr -d '[ -~]\0' | wc -c) != 0
then
	cat <<\EOF
Error: Attempt to add a non-ASCII file name.

This can cause problems if you want to work with people on other platforms.

To be portable it is advisable to rename the file.

If you know what you are doing you can disable this check using:

  git config hooks.allownonascii true
EOF
	exit 1
fi

# If there are whitespace errors, print the offending file names and fail.
exec git diff-index --check --cached $against --

```

## .git/hooks/pre-merge-commit.sample

```text
#!/bin/sh
#
# An example hook script to verify what is about to be committed.
# Called by "git merge" with no arguments.  The hook should
# exit with non-zero status after issuing an appropriate message to
# stderr if it wants to stop the merge commit.
#
# To enable this hook, rename this file to "pre-merge-commit".

. git-sh-setup
test -x "$GIT_DIR/hooks/pre-commit" &&
        exec "$GIT_DIR/hooks/pre-commit"
:

```

## .git/hooks/pre-push.sample

```text
#!/bin/sh

# An example hook script to verify what is about to be pushed.  Called by "git
# push" after it has checked the remote status, but before anything has been
# pushed.  If this script exits with a non-zero status nothing will be pushed.
#
# This hook is called with the following parameters:
#
# $1 -- Name of the remote to which the push is being done
# $2 -- URL to which the push is being done
#
# If pushing without using a named remote those arguments will be equal.
#
# Information about the commits which are being pushed is supplied as lines to
# the standard input in the form:
#
#   <local ref> <local oid> <remote ref> <remote oid>
#
# This sample shows how to prevent push of commits where the log message starts
# with "WIP" (work in progress).

remote="$1"
url="$2"

zero=$(git hash-object --stdin </dev/null | tr '[0-9a-f]' '0')

while read local_ref local_oid remote_ref remote_oid
do
	if test "$local_oid" = "$zero"
	then
		# Handle delete
		:
	else
		if test "$remote_oid" = "$zero"
		then
			# New branch, examine all commits
			range="$local_oid"
		else
			# Update to existing branch, examine new commits
			range="$remote_oid..$local_oid"
		fi

		# Check for WIP commit
		commit=$(git rev-list -n 1 --grep '^WIP' "$range")
		if test -n "$commit"
		then
			echo >&2 "Found WIP commit in $local_ref, not pushing"
			exit 1
		fi
	fi
done

exit 0

```

## .git/hooks/pre-rebase.sample

```text
#!/bin/sh
#
# Copyright (c) 2006, 2008 Junio C Hamano
#
# The "pre-rebase" hook is run just before "git rebase" starts doing
# its job, and can prevent the command from running by exiting with
# non-zero status.
#
# The hook is called with the following parameters:
#
# $1 -- the upstream the series was forked from.
# $2 -- the branch being rebased (or empty when rebasing the current branch).
#
# This sample shows how to prevent topic branches that are already
# merged to 'next' branch from getting rebased, because allowing it
# would result in rebasing already published history.

publish=next
basebranch="$1"
if test "$#" = 2
then
	topic="refs/heads/$2"
else
	topic=`git symbolic-ref HEAD` ||
	exit 0 ;# we do not interrupt rebasing detached HEAD
fi

case "$topic" in
refs/heads/??/*)
	;;
*)
	exit 0 ;# we do not interrupt others.
	;;
esac

# Now we are dealing with a topic branch being rebased
# on top of master.  Is it OK to rebase it?

# Does the topic really exist?
git show-ref -q "$topic" || {
	echo >&2 "No such branch $topic"
	exit 1
}

# Is topic fully merged to master?
not_in_master=`git rev-list --pretty=oneline ^master "$topic"`
if test -z "$not_in_master"
then
	echo >&2 "$topic is fully merged to master; better remove it."
	exit 1 ;# we could allow it, but there is no point.
fi

# Is topic ever merged to next?  If so you should not be rebasing it.
only_next_1=`git rev-list ^master "^$topic" ${publish} | sort`
only_next_2=`git rev-list ^master           ${publish} | sort`
if test "$only_next_1" = "$only_next_2"
then
	not_in_topic=`git rev-list "^$topic" master`
	if test -z "$not_in_topic"
	then
		echo >&2 "$topic is already up to date with master"
		exit 1 ;# we could allow it, but there is no point.
	else
		exit 0
	fi
else
	not_in_next=`git rev-list --pretty=oneline ^${publish} "$topic"`
	/usr/bin/perl -e '
		my $topic = $ARGV[0];
		my $msg = "* $topic has commits already merged to public branch:\n";
		my (%not_in_next) = map {
			/^([0-9a-f]+) /;
			($1 => 1);
		} split(/\n/, $ARGV[1]);
		for my $elem (map {
				/^([0-9a-f]+) (.*)$/;
				[$1 => $2];
			} split(/\n/, $ARGV[2])) {
			if (!exists $not_in_next{$elem->[0]}) {
				if ($msg) {
					print STDERR $msg;
					undef $msg;
				}
				print STDERR " $elem->[1]\n";
			}
		}
	' "$topic" "$not_in_next" "$not_in_master"
	exit 1
fi

<<\DOC_END

This sample hook safeguards topic branches that have been
published from being rewound.

The workflow assumed here is:

 * Once a topic branch forks from "master", "master" is never
   merged into it again (either directly or indirectly).

 * Once a topic branch is fully cooked and merged into "master",
   it is deleted.  If you need to build on top of it to correct
   earlier mistakes, a new topic branch is created by forking at
   the tip of the "master".  This is not strictly necessary, but
   it makes it easier to keep your history simple.

 * Whenever you need to test or publish your changes to topic
   branches, merge them into "next" branch.

The script, being an example, hardcodes the publish branch name
to be "next", but it is trivial to make it configurable via
$GIT_DIR/config mechanism.

With this workflow, you would want to know:

(1) ... if a topic branch has ever been merged to "next".  Young
    topic branches can have stupid mistakes you would rather
    clean up before publishing, and things that have not been
    merged into other branches can be easily rebased without
    affecting other people.  But once it is published, you would
    not want to rewind it.

(2) ... if a topic branch has been fully merged to "master".
    Then you can delete it.  More importantly, you should not
    build on top of it -- other people may already want to
    change things related to the topic as patches against your
    "master", so if you need further changes, it is better to
    fork the topic (perhaps with the same name) afresh from the
    tip of "master".

Let's look at this example:

		   o---o---o---o---o---o---o---o---o---o "next"
		  /       /           /           /
		 /   a---a---b A     /           /
		/   /               /           /
	       /   /   c---c---c---c B         /
	      /   /   /             \         /
	     /   /   /   b---b C     \       /
	    /   /   /   /             \     /
    ---o---o---o---o---o---o---o---o---o---o---o "master"


A, B and C are topic branches.

 * A has one fix since it was merged up to "next".

 * B has finished.  It has been fully merged up to "master" and "next",
   and is ready to be deleted.

 * C has not merged to "next" at all.

We would want to allow C to be rebased, refuse A, and encourage
B to be deleted.

To compute (1):

	git rev-list ^master ^topic next
	git rev-list ^master        next

	if these match, topic has not merged in next at all.

To compute (2):

	git rev-list master..topic

	if this is empty, it is fully merged to "master".

DOC_END

```

## .git/hooks/pre-receive.sample

```text
#!/bin/sh
#
# An example hook script to make use of push options.
# The example simply echoes all push options that start with 'echoback='
# and rejects all pushes when the "reject" push option is used.
#
# To enable this hook, rename this file to "pre-receive".

if test -n "$GIT_PUSH_OPTION_COUNT"
then
	i=0
	while test "$i" -lt "$GIT_PUSH_OPTION_COUNT"
	do
		eval "value=\$GIT_PUSH_OPTION_$i"
		case "$value" in
		echoback=*)
			echo "echo from the pre-receive-hook: ${value#*=}" >&2
			;;
		reject)
			exit 1
		esac
		i=$((i + 1))
	done
fi

```

## .git/hooks/prepare-commit-msg.sample

```text
#!/bin/sh
#
# An example hook script to prepare the commit log message.
# Called by "git commit" with the name of the file that has the
# commit message, followed by the description of the commit
# message's source.  The hook's purpose is to edit the commit
# message file.  If the hook fails with a non-zero status,
# the commit is aborted.
#
# To enable this hook, rename this file to "prepare-commit-msg".

# This hook includes three examples. The first one removes the
# "# Please enter the commit message..." help message.
#
# The second includes the output of "git diff --name-status -r"
# into the message, just before the "git status" output.  It is
# commented because it doesn't cope with --amend or with squashed
# commits.
#
# The third example adds a Signed-off-by line to the message, that can
# still be edited.  This is rarely a good idea.

COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2
SHA1=$3

/usr/bin/perl -i.bak -ne 'print unless(m/^. Please enter the commit message/..m/^#$/)' "$COMMIT_MSG_FILE"

# case "$COMMIT_SOURCE,$SHA1" in
#  ,|template,)
#    /usr/bin/perl -i.bak -pe '
#       print "\n" . `git diff --cached --name-status -r`
# 	 if /^#/ && $first++ == 0' "$COMMIT_MSG_FILE" ;;
#  *) ;;
# esac

# SOB=$(git var GIT_COMMITTER_IDENT | sed -n 's/^\(.*>\).*$/Signed-off-by: \1/p')
# git interpret-trailers --in-place --trailer "$SOB" "$COMMIT_MSG_FILE"
# if test -z "$COMMIT_SOURCE"
# then
#   /usr/bin/perl -i.bak -pe 'print "\n" if !$first_line++' "$COMMIT_MSG_FILE"
# fi

```

## .git/hooks/push-to-checkout.sample

```text
#!/bin/sh

# An example hook script to update a checked-out tree on a git push.
#
# This hook is invoked by git-receive-pack(1) when it reacts to git
# push and updates reference(s) in its repository, and when the push
# tries to update the branch that is currently checked out and the
# receive.denyCurrentBranch configuration variable is set to
# updateInstead.
#
# By default, such a push is refused if the working tree and the index
# of the remote repository has any difference from the currently
# checked out commit; when both the working tree and the index match
# the current commit, they are updated to match the newly pushed tip
# of the branch. This hook is to be used to override the default
# behaviour; however the code below reimplements the default behaviour
# as a starting point for convenient modification.
#
# The hook receives the commit with which the tip of the current
# branch is going to be updated:
commit=$1

# It can exit with a non-zero status to refuse the push (when it does
# so, it must not modify the index or the working tree).
die () {
	echo >&2 "$*"
	exit 1
}

# Or it can make any necessary changes to the working tree and to the
# index to bring them to the desired state when the tip of the current
# branch is updated to the new commit, and exit with a zero status.
#
# For example, the hook can simply run git read-tree -u -m HEAD "$1"
# in order to emulate git fetch that is run in the reverse direction
# with git push, as the two-tree form of git read-tree -u -m is
# essentially the same as git switch or git checkout that switches
# branches while keeping the local changes in the working tree that do
# not interfere with the difference between the branches.

# The below is a more-or-less exact translation to shell of the C code
# for the default behaviour for git's push-to-checkout hook defined in
# the push_to_deploy() function in builtin/receive-pack.c.
#
# Note that the hook will be executed from the repository directory,
# not from the working tree, so if you want to perform operations on
# the working tree, you will have to adapt your code accordingly, e.g.
# by adding "cd .." or using relative paths.

if ! git update-index -q --ignore-submodules --refresh
then
	die "Up-to-date check failed"
fi

if ! git diff-files --quiet --ignore-submodules --
then
	die "Working directory has unstaged changes"
fi

# This is a rough translation of:
#
#   head_has_history() ? "HEAD" : EMPTY_TREE_SHA1_HEX
if git cat-file -e HEAD 2>/dev/null
then
	head=HEAD
else
	head=$(git hash-object -t tree --stdin </dev/null)
fi

if ! git diff-index --quiet --cached --ignore-submodules $head --
then
	die "Working directory has staged changes"
fi

if ! git read-tree -u -m "$commit"
then
	die "Could not update working tree to new HEAD"
fi

```

## .git/hooks/sendemail-validate.sample

```text
#!/bin/sh

# An example hook script to validate a patch (and/or patch series) before
# sending it via email.
#
# The hook should exit with non-zero status after issuing an appropriate
# message if it wants to prevent the email(s) from being sent.
#
# To enable this hook, rename this file to "sendemail-validate".
#
# By default, it will only check that the patch(es) can be applied on top of
# the default upstream branch without conflicts in a secondary worktree. After
# validation (successful or not) of the last patch of a series, the worktree
# will be deleted.
#
# The following config variables can be set to change the default remote and
# remote ref that are used to apply the patches against:
#
#   sendemail.validateRemote (default: origin)
#   sendemail.validateRemoteRef (default: HEAD)
#
# Replace the TODO placeholders with appropriate checks according to your
# needs.

validate_cover_letter () {
	file="$1"
	# TODO: Replace with appropriate checks (e.g. spell checking).
	true
}

validate_patch () {
	file="$1"
	# Ensure that the patch applies without conflicts.
	git am -3 "$file" || return
	# TODO: Replace with appropriate checks for this patch
	# (e.g. checkpatch.pl).
	true
}

validate_series () {
	# TODO: Replace with appropriate checks for the whole series
	# (e.g. quick build, coding style checks, etc.).
	true
}

# main -------------------------------------------------------------------------

if test "$GIT_SENDEMAIL_FILE_COUNTER" = 1
then
	remote=$(git config --default origin --get sendemail.validateRemote) &&
	ref=$(git config --default HEAD --get sendemail.validateRemoteRef) &&
	worktree=$(mktemp --tmpdir -d sendemail-validate.XXXXXXX) &&
	git worktree add -fd --checkout "$worktree" "refs/remotes/$remote/$ref" &&
	git config --replace-all sendemail.validateWorktree "$worktree"
else
	worktree=$(git config --get sendemail.validateWorktree)
fi || {
	echo "sendemail-validate: error: failed to prepare worktree" >&2
	exit 1
}

unset GIT_DIR GIT_WORK_TREE
cd "$worktree" &&

if grep -q "^diff --git " "$1"
then
	validate_patch "$1"
else
	validate_cover_letter "$1"
fi &&

if test "$GIT_SENDEMAIL_FILE_COUNTER" = "$GIT_SENDEMAIL_FILE_TOTAL"
then
	git config --unset-all sendemail.validateWorktree &&
	trap 'git worktree remove -ff "$worktree"' EXIT &&
	validate_series
fi

```

## .git/hooks/update.sample

```text
#!/bin/sh
#
# An example hook script to block unannotated tags from entering.
# Called by "git receive-pack" with arguments: refname sha1-old sha1-new
#
# To enable this hook, rename this file to "update".
#
# Config
# ------
# hooks.allowunannotated
#   This boolean sets whether unannotated tags will be allowed into the
#   repository.  By default they won't be.
# hooks.allowdeletetag
#   This boolean sets whether deleting tags will be allowed in the
#   repository.  By default they won't be.
# hooks.allowmodifytag
#   This boolean sets whether a tag may be modified after creation. By default
#   it won't be.
# hooks.allowdeletebranch
#   This boolean sets whether deleting branches will be allowed in the
#   repository.  By default they won't be.
# hooks.denycreatebranch
#   This boolean sets whether remotely creating branches will be denied
#   in the repository.  By default this is allowed.
#

# --- Command line
refname="$1"
oldrev="$2"
newrev="$3"

# --- Safety check
if [ -z "$GIT_DIR" ]; then
	echo "Don't run this script from the command line." >&2
	echo " (if you want, you could supply GIT_DIR then run" >&2
	echo "  $0 <ref> <oldrev> <newrev>)" >&2
	exit 1
fi

if [ -z "$refname" -o -z "$oldrev" -o -z "$newrev" ]; then
	echo "usage: $0 <ref> <oldrev> <newrev>" >&2
	exit 1
fi

# --- Config
allowunannotated=$(git config --type=bool hooks.allowunannotated)
allowdeletebranch=$(git config --type=bool hooks.allowdeletebranch)
denycreatebranch=$(git config --type=bool hooks.denycreatebranch)
allowdeletetag=$(git config --type=bool hooks.allowdeletetag)
allowmodifytag=$(git config --type=bool hooks.allowmodifytag)

# check for no description
projectdesc=$(sed -e '1q' "$GIT_DIR/description")
case "$projectdesc" in
"Unnamed repository"* | "")
	echo "*** Project description file hasn't been set" >&2
	exit 1
	;;
esac

# --- Check types
# if $newrev is 0000...0000, it's a commit to delete a ref.
zero=$(git hash-object --stdin </dev/null | tr '[0-9a-f]' '0')
if [ "$newrev" = "$zero" ]; then
	newrev_type=delete
else
	newrev_type=$(git cat-file -t $newrev)
fi

case "$refname","$newrev_type" in
	refs/tags/*,commit)
		# un-annotated tag
		short_refname=${refname##refs/tags/}
		if [ "$allowunannotated" != "true" ]; then
			echo "*** The un-annotated tag, $short_refname, is not allowed in this repository" >&2
			echo "*** Use 'git tag [ -a | -s ]' for tags you want to propagate." >&2
			exit 1
		fi
		;;
	refs/tags/*,delete)
		# delete tag
		if [ "$allowdeletetag" != "true" ]; then
			echo "*** Deleting a tag is not allowed in this repository" >&2
			exit 1
		fi
		;;
	refs/tags/*,tag)
		# annotated tag
		if [ "$allowmodifytag" != "true" ] && git rev-parse $refname > /dev/null 2>&1
		then
			echo "*** Tag '$refname' already exists." >&2
			echo "*** Modifying a tag is not allowed in this repository." >&2
			exit 1
		fi
		;;
	refs/heads/*,commit)
		# branch
		if [ "$oldrev" = "$zero" -a "$denycreatebranch" = "true" ]; then
			echo "*** Creating a branch is not allowed in this repository" >&2
			exit 1
		fi
		;;
	refs/heads/*,delete)
		# delete branch
		if [ "$allowdeletebranch" != "true" ]; then
			echo "*** Deleting a branch is not allowed in this repository" >&2
			exit 1
		fi
		;;
	refs/remotes/*,commit)
		# tracking branch
		;;
	refs/remotes/*,delete)
		# delete tracking branch
		if [ "$allowdeletebranch" != "true" ]; then
			echo "*** Deleting a tracking branch is not allowed in this repository" >&2
			exit 1
		fi
		;;
	*)
		# Anything else (is there anything else?)
		echo "*** Update hook: unknown type of update to ref $refname of type $newrev_type" >&2
		exit 1
		;;
esac

# --- Finished
exit 0

```

## .git/index

> Fichier binaire non inclus (4441 octets)

## .git/info/exclude

```text
# git ls-files --others --exclude-from=.git/info/exclude
# Lines that start with '#' are comments.
# For a project mostly in C, the following would be a good set of
# exclude patterns (uncomment them if you want to use them):
# *.[oa]
# *~

```

## .git/logs/HEAD

```text
0000000000000000000000000000000000000000 9dcf46457a2239e54fbe4a9f00f2826996522e04 sicDANGBE <dansoug@gmail.com> 1766448047 +0100	commit (initial): feat: initial gRPC fibonacci benchmark with Go, Python and Node
9dcf46457a2239e54fbe4a9f00f2826996522e04 c9a83108423860e7c4e872bdb748dfb3d2f60be9 sicDANGBE <dansoug@gmail.com> 1766677162 +0100	commit: v2 orchestration + ajout rust
c9a83108423860e7c4e872bdb748dfb3d2f60be9 8913571fad3aed041bd6f5a6fbb57446cf28ec84 sicDANGBE <dansoug@gmail.com> 1766681898 +0100	commit: docs: ajouter .github/copilot-instructions.md ‚Äî guide agent Copilot pour le benchmark
8913571fad3aed041bd6f5a6fbb57446cf28ec84 cb01ac9c40ec00fd41ff670bbc86735946ff34a9 sicDANGBE <dansoug@gmail.com> 1766682904 +0100	commit: ui: am√©liorer affichage des cartes worker et indicateur WebSocket
cb01ac9c40ec00fd41ff670bbc86735946ff34a9 84dc0203f93f52853e23371eb3c46d49449ccf3e sicDANGBE <dansoug@gmail.com> 1766682917 +0100	commit: orchestrator: robustifier consommation des streams RabbitMQ (QoS, ack manuel, offset)
84dc0203f93f52853e23371eb3c46d49449ccf3e 5a8af5557d17aa3d37ee1e7841508587a0d49f19 sicDANGBE <dansoug@gmail.com> 1766682924 +0100	commit: worker-go: publier metrics et enrichir metadata dans handler fibonacci (√©viter QueueDeclare redondant)

```

## .git/logs/refs/heads/master

```text
0000000000000000000000000000000000000000 9dcf46457a2239e54fbe4a9f00f2826996522e04 sicDANGBE <dansoug@gmail.com> 1766448047 +0100	commit (initial): feat: initial gRPC fibonacci benchmark with Go, Python and Node
9dcf46457a2239e54fbe4a9f00f2826996522e04 c9a83108423860e7c4e872bdb748dfb3d2f60be9 sicDANGBE <dansoug@gmail.com> 1766677162 +0100	commit: v2 orchestration + ajout rust
c9a83108423860e7c4e872bdb748dfb3d2f60be9 8913571fad3aed041bd6f5a6fbb57446cf28ec84 sicDANGBE <dansoug@gmail.com> 1766681898 +0100	commit: docs: ajouter .github/copilot-instructions.md ‚Äî guide agent Copilot pour le benchmark
8913571fad3aed041bd6f5a6fbb57446cf28ec84 cb01ac9c40ec00fd41ff670bbc86735946ff34a9 sicDANGBE <dansoug@gmail.com> 1766682904 +0100	commit: ui: am√©liorer affichage des cartes worker et indicateur WebSocket
cb01ac9c40ec00fd41ff670bbc86735946ff34a9 84dc0203f93f52853e23371eb3c46d49449ccf3e sicDANGBE <dansoug@gmail.com> 1766682917 +0100	commit: orchestrator: robustifier consommation des streams RabbitMQ (QoS, ack manuel, offset)
84dc0203f93f52853e23371eb3c46d49449ccf3e 5a8af5557d17aa3d37ee1e7841508587a0d49f19 sicDANGBE <dansoug@gmail.com> 1766682924 +0100	commit: worker-go: publier metrics et enrichir metadata dans handler fibonacci (√©viter QueueDeclare redondant)

```

## .git/logs/refs/remotes/origin/master

```text
0000000000000000000000000000000000000000 9dcf46457a2239e54fbe4a9f00f2826996522e04 sicDANGBE <dansoug@gmail.com> 1766448056 +0100	update by push
9dcf46457a2239e54fbe4a9f00f2826996522e04 c9a83108423860e7c4e872bdb748dfb3d2f60be9 sicDANGBE <dansoug@gmail.com> 1766677171 +0100	update by push
c9a83108423860e7c4e872bdb748dfb3d2f60be9 8913571fad3aed041bd6f5a6fbb57446cf28ec84 sicDANGBE <dansoug@gmail.com> 1766681910 +0100	update by push
8913571fad3aed041bd6f5a6fbb57446cf28ec84 5a8af5557d17aa3d37ee1e7841508587a0d49f19 sicDANGBE <dansoug@gmail.com> 1766682948 +0100	update by push

```

## .git/objects/07/ae975db35e66033dec0b124a5beb703025248f

> Fichier binaire non inclus (407 octets)

## .git/objects/0c/4fc5ef3fa24178bd7db7c90efb69907387eee0

> Fichier binaire non inclus (83 octets)

## .git/objects/0c/636d86bffac4f802278a4a7178d9a27193f0ec

> Fichier binaire non inclus (610 octets)

## .git/objects/0e/ee988462b970f9b69b1449ac6c54ca6a9a5754

> Fichier binaire non inclus (152 octets)

## .git/objects/11/c0a428ae63867970c057e4da81ea673a9b10a4

> Fichier binaire non inclus (52 octets)

## .git/objects/12/a122ea7678e77dcb0ffc86aeca9c7d69bd816d

> Fichier binaire non inclus (848 octets)

## .git/objects/12/d88b2843b5c4a93c94343f51967ec3b4538f6d

> Fichier binaire non inclus (203 octets)

## .git/objects/16/f77612eccd94213863c5ad7ebd2e9fe8dd009b

> Fichier binaire non inclus (236 octets)

## .git/objects/1a/fb529027bfd3b244eb797399e6c0966eecbeb5

> Fichier binaire non inclus (117 octets)

## .git/objects/1d/17b6c4d2312673094dc22ccfeac6262f182d81

> Fichier binaire non inclus (48 octets)

## .git/objects/1e/9e3512be43ceac03f749b3e1bd81a4e0899e7f

> Fichier binaire non inclus (956 octets)

## .git/objects/20/708f6a2ef877be2ccc2101160dbce8791b4604

```text
xUè¡Jƒ0Ü=Á)ÜzX”¢{POYï"ª]¬äàH…¶—¥ôê§`©œ—≥m
‚)ì˘gÊˇ˛cçGXﬂﬁù=Òb∂'4˜ÎÙ˙Ü˙Z7Ñº¸eìs»ÑµÑ◊XmAD]•©Ú§h•|;+5.àµ'‰Z≠†j¡:MP?AÅRuIX±èbÊ;#”8óNªœCoÜﬁâ†—L'$V
ˆ3f§ôK†ÕÏZŒûq_Õ”ë0f*±”wb,ˇı˛,'CJ)∞ÇÛGv»ã‰,üZ£PvË•˛√‘cÙ^+['¨‚Õlµ§ôπ≤Fhì⁄∆,lªÅè$N%Wê,RÚ˘[ÇÄk
```

## .git/objects/22/1181d2632a36a783df61e3c0345b3ede421068

> Fichier binaire non inclus (235 octets)

## .git/objects/22/f7a387d14bb1ff975b8f8a909c2760552ab6b2

```text
xùTÔO€0›Á˛ßnR	€Ù√–§MÇ}”> N|-ﬁ;ÚZÑ˙øÔú§•§Eä‚8wÔΩª{N©m	„Ω√Ò\4÷òFSeƒFäÄﬂÉ—gï5´Ä2áá–E>Ä2RU"X' mk4°òa¯¨1-?›ïŸhÓôoQF˘Qõ´¶0ÅÆ·äJÔœîÖ√Fã
≥Q9c%;‡|¥Èç(ú–›Nè¸*àø“ŒŸ%øÊ◊cﬁ,Æ	Ú*·=˘vÃÀè∆W+—K@Ì±Ø˛Uûø‰uÇWÚˇ[ÏZ”sÇ€JzËÂ`9<?–ü÷˝FóZõÕ€ÂˆD5}xqòmKa´ŒÙVêîv”A≤w›¢PryÛ8ÙYD…<á!:s4ÿ∞ÍM)ïC2do≠l$’›Jj"¢P%;.zoÙM‘H€√©∆®ÄµgŸ¸ä>®È=+1Ã4lhf^{œ˘Óg£ë‰∏ÖÜ“:IY›£è:‰ÑQu h"yc∏&WtX‹óãÛ≥‘í∂ûdòcR≠ÆìNQí≈*´áìuLÁa÷Å÷VZ-!‡ÇV:"€Á|8Y∑X3ãbÜE∞?ö›©òÂÀ„›Ûrãwô¡LKmçÌX∫6l”(Y¯X˙‡îôe|g|¯î‚xó*|d‹¨wŒˆ‡ñÓUOßQSW∑Œp≤ëﬂõ&Ÿ¨Tïëß∑JÀuÔpè·B’hc»≤N&t:∑~µΩ£∆÷|F˘ÏqŒ	a˘ëÎê˝
```

## .git/objects/24/3e0c81c7ad72a6d3e8acb2dec9be096a1193da

> Fichier binaire non inclus (590 octets)

## .git/objects/26/1363925bedc249329ef083cf913018d63c62f5

> Fichier binaire non inclus (565 octets)

## .git/objects/33/c0d87c6d62f2f9d9a4e40f8f25ed90e4056ba3

> Fichier binaire non inclus (110 octets)

## .git/objects/34/2ec37e7c6d81a03252d768ba5b12b3f726e491

> Fichier binaire non inclus (472 octets)

## .git/objects/37/ae47ae280e513c8aa2129bd8132aefba977851

> Fichier binaire non inclus (5707 octets)

## .git/objects/3d/2251a550284deb468cc7e650b8b3aee7127908

> Fichier binaire non inclus (604 octets)

## .git/objects/3e/2093992b9bec6d3f7c01739e68ac39c300d3d1

> Fichier binaire non inclus (125 octets)

## .git/objects/3e/20ed345bffc4e3750b64a9529fa7d180030b91

```text
xÖR€J1ıµ˚Cˇ@ê
æ|∞† ZA§ùn‚6ÓÊB2ãJÈøõÎ≤ª(Ê%ìô3sÊrËÙVó◊Î⁄÷GÓ»"i[UÙm8ºh€r˚ƒÛB+A_ú™≈Ê“Ò)°ÿ8≠nñÇ-˜’‚U”á©”bó”rŒk&Ö⁄¢kG£√3åœÕôg∏2y‘.—›°b∑ˇ4*Ï˜LhiM ≠ÆÜ¡Cê\@ÌêB«#Zî@¢yM{Ω˘fnﬂ±ÊßsQo"j$ØXË˙éF
SzÆ1”~F€≥¥‚≈‘çø-¯7RæQåE°~˝π`äAÌ'dH¡°πHôã∏“îfbeYØ£!?Íz≈
```

## .git/objects/3f/a8c8c5920b5fd88bbfd47984588cab3f7465c5

> Fichier binaire non inclus (370 octets)

## .git/objects/42/4bb4fbbc66cf6af15cf54a2ff469815c3faddf

> Fichier binaire non inclus (119 octets)

## .git/objects/47/1c30f56909f444ef7f38a7049ab5233f1de6fc

> Fichier binaire non inclus (114 octets)

## .git/objects/49/6c619e5bce5b8e4b24d430a73bbd67e2f4f31e

> Fichier binaire non inclus (15640202 octets)

## .git/objects/49/a6f4e3b2dce8aee8a067bf9e18b853b66fd19d

> Fichier binaire non inclus (81 octets)

## .git/objects/4b/37f9a11f378cc7e66ce3371d0508022174dd30

```text
xMè¡N√0D9˚+F9µ–Q/î*!"µë8"«Ÿ´±◊Z;@UıÉ¯~SQƒﬁfvvû∂∏≈Ï˙ˆ"Ó}“ò£¬âgE•T–fß{Bﬁô,ÀãzΩ^.öU˝Ñ;<Ã+9ÎÒPØvéº!A–ÇÅ‡Ÿ°·∏≥Íæ>∑÷[t⁄GÙ|ïm≈!YˆYΩúQôøµ-_æ≥ÏHJÎâ◊C⁄¢ÇRë‰Õ¬Ω±ôtP»#¡‡Y€‘&iIò,]H˚)Ñ“(ô69Ÿ€Á¶i•éJ9äÒÁ≥Sá„üÒ/˘[ìXﬂ„|1«Mn¯Åòg‘
```

## .git/objects/4c/aeed4572829ae912ba1a0d88c87af1bc82e713

> Fichier binaire non inclus (1202 octets)

## .git/objects/51/892cf7e7cd672d3b5d6d471bbb562d04a8e164

> Fichier binaire non inclus (75 octets)

## .git/objects/5a/50d4c26668c2b5b465756529b5868e7f109e11

> Fichier binaire non inclus (1105 octets)

## .git/objects/5a/8af5557d17aa3d37ee1e7841508587a0d49f19

```text
xïOKN≈ u‹U‹°∆h((5∆®y∆ôâK∏\.-±ÖJuMÆ√ç…sŒNNŒóÚ∫∆
™µ0√‡åë=kÁ{>∏`Ö‘É±H“Ñ`¨Â`ï∞≤;a·T¡jOB
F“K•Üûù"mºµâÇ‚˜:Á[§„”€ÎÛ‹{L[ﬁß«i≈∏‹R^†å1VéR√µËÖË€ÊU˛∑±˚ ÂÉÀÕîÔ‡¥ª%∂àïkâ¥W‡‘–ˇ8ÙXŒk`∆‰ó¶—ÂÑD.æ?„πˇ}ÁùèLK˚Ö}NÕQØ∫_∏0f	
```

## .git/objects/5a/e095b7d3a678a69e6808cb3856ccd3ded0bdf5

> Fichier binaire non inclus (1301 octets)

## .git/objects/5d/175e2933d1e59016f42d4e2f2441db2fef311d

> Fichier binaire non inclus (117 octets)

## .git/objects/5d/ab255ae1904e7a27735e3c9641a108451695d3

```text
x}ê¡nÑ0D{mæ"«ˆÄ«ËÁQ$≥V€øoñ=¥e•ﬁ,k<Û∆˝ÃΩ$€º|Ú∏œ^N±Á‰Ü!
X"òPàÏ/{Ã^Ê0{<ªÄsP!/Éº"¥-Ë›õx˝•π©‰7y’P[¿JÉ—¶Açd®n©´,çË BOùëJ…ò∆í5l'èıkΩ{P	˙Gµ˘€E˙YˆƒÓ”íycı®Âñ∏™£Nπ]=0ç¡∫L∂"ám=}ÿŒ∫Ópv>l˚}∫ÜJÈøƒÔ‚q◊r
```

## .git/objects/5e/1af68ac259b3f3f34291e2803d602d13ce9fd5

```text
xMO=KAµﬁ_ÒHsIu®BH#XXG∞î˝ 1x∑≥ŒÓâ!‰ø;{p¡bä7Ôk∆çÏ¸tWŒ©⁄_∞…¬ï7{c≤ı_vàPŒ+Ï{º•@ﬂsDõ1bú,UR)"Ãò8Ã v'rú¨˜‘!ÿTö#p)]vù·\â˛\É¥˛fÍ≥kGî(?‰#^¨Hs^ Ÿ„C;ﬂ˘Xµ€◊)◊ÛÎ,Z¥]÷Gíw{s5fä•¥W!.◊€‚ürâ.U(Xı‹´ˇóöc'
```

## .git/objects/5f/7131495b2d34d83cb08d08a31f739f2fa24949

```text
xùî=n‹0ÖS˚SÓ±LÍoe)Ç©í*ÈåÖ°•hô∂ƒ!H Y•Ú!“§KiùC7ÒIL…Zƒ±#òH«bæ7√y3≥´q4äì7g!	ìcáí«$ßt_±’å‰-Eh+¶•Òä3{Œ˜
µörÅ¯$ÿ•‡⁄Ä’ÿﬁΩÅUqSH¢∂öØs†…kËJVÿ°ó\;•á;ÛD!Üﬂ≥∫5P¥s
•–9xN÷Ú%Ô$á3˜yÜçB√ÉÆ©aeƒ˛.=›¯∞·*<˘àÏöÎ1„Lg$Ù°£âÆ0h∞ú…(M|»¯@ö∂9ê4ÛJöLhST8≥!IºæõN¨⁄ùòN≤@Ì˛(xuzÛ?Ø¥˙çl◊N„MÍ”≤”-H,˘K≥BÍ’rJf!KæÆ‹0OìíE^È«{∂75VçƒÀ:Í∆ÕiXú˚?>¶^R7r™≥ó(_ˆ!ä#ü>R7{≥ƒdÑÍÊRöΩ¬˜≤∞%|)Ùuâﬂ%XÆ!áûª≥Cc∑≥”ΩXôı[†tC	 ≥‹ö•ãÙæe≠€`w*∑zÂ˝ÌØB¢ÏaÛ‹ﬂ˛Ñ°∑Ö‚`‹qpyñîæÕÖŸôRéê®Î°œˇ·Ÿ—ØZõi
```

## .git/objects/5f/b088a58ba421fe022337f7079466ed09d22bb1

```text
xUé]k1E}ŒØ|œUƒ"¯ äPäÆl)Eƒádw∂]ù0I¸˝›/ê>ŒeÓ=«:∂0[.&˚"?Ä§W/ŸÎ\Áõ;¡Êlj\E¢äœ#ˇ¶™@Î;Î“î?∑úÆË^8≤Muw(ıïÔª∑0¡ %Ô’6?ùá?Ñ,√æÅC⁄ë€∞ßîFæy‡∂$!G&êRΩ·‡µr&RàOä0«qIÎZ¯∂µˇ	`ló)‚8âucYw‰º=Ï‡2ÕûŸÙ˙∫\Z
```

## .git/objects/61/9c750eb7e119698413b3f6346eb8c14f603684

> Fichier binaire non inclus (119 octets)

## .git/objects/64/72bab105a7266bf1324976307cfdc40c9fba6d

> Fichier binaire non inclus (1067 octets)

## .git/objects/67/057bc258203607838b49236920023738fec519

> Fichier binaire non inclus (93 octets)

## .git/objects/67/733188b7153a7ce368444a9509968142172eea

```text
xÖRMo‘0‰ˇäß\»™ê4°õ/i%8 @jÅnUıX˘„%kmbß∂”"Pˇ;vvª8ÙËÁy3û≥A3(ã¸ÕD˘ûˆ#ïä9N⁄8HHwíÈ˜⁄Zg®√ŸdR94äùd¸
‰ﬂ¶6´mLVÑt≥‚ãV≤Ç?$ 2»S¯¶§ìtêñ:©àæŒÓê›hæGG¢ù?∂≤Èw|Úó…äDΩ?O∑≥ÚßÖ™¯ü
ax˚YıR!$[ òtW◊p?Næº†Á¢„√tªΩ"⁄¶_–°zL‚OW◊?Ô˝8ˆŸ¡h≥Å8Ø?≠˘A∏l≥å°‚ª˚©@”>XQÆªﬁ‘ÇWùhl]ı¥Ï¨X˜ø≠0]”◊U#:[ıÛ¶HÛ≤NÛ4/⁄uY™ëöΩO.z>8ÙÓèn™∏7GÅÈô¡ ◊J·Øê!*8¥Ö$
U,ŒN¶µ	9òí£≠w!ÕcêR∏A7OÄ.É–à Ö^,öGˇB7$2/ù,ÿ≠û˝ÔHÇÿë*2K9q[ü◊Á>√gÚÀ„«›
```

## .git/objects/69/debbb92bbcb24a4513e4cf15034b3f40c2f338

> Fichier binaire non inclus (52 octets)

## .git/objects/6c/cc48a05421e0f3e40a0cfd4bc069663295c1c7

> Fichier binaire non inclus (32 octets)

## .git/objects/6d/6408a5bb8411448f2ee0f3e05078f8e13ff684

> Fichier binaire non inclus (1673 octets)

## .git/objects/6f/9367e4ddeff2ceee6e1b6aa24bbf00f029cb75

> Fichier binaire non inclus (1609 octets)

## .git/objects/71/d1529cacb1c8b66dff4899a11595712c0ff2ca

> Fichier binaire non inclus (1445 octets)

## .git/objects/73/a6ec090865db226e94fa131701d00746aaa7a9

> Fichier binaire non inclus (407 octets)

## .git/objects/75/cdea54465f53768e9edc7c31fa5a226602af79

> Fichier binaire non inclus (82 octets)

## .git/objects/75/cfbdfbadfe07187291e5446058a8c6349db908

> Fichier binaire non inclus (203 octets)

## .git/objects/78/75ca23928e02658b93e6524a8c99458e7c39e6

```text
x+)JMU02∂`040031Q–KÃ,“+…œÕaxƒ«i2o‚ãSˇX>úUΩ¥≈s≥Îƒ®"ó¸‰Ï‘¢¥ÃúTΩ√u5πçåÇÆgÏä⁄¸]Ì…D PHŒMaX∫ÛXÈ=«“å⁄˙ˇÅnmeˇ-ÿ5#=_/7?ÖA≈éßÒ¯⁄¢eó_¨ŸtÔ‰>Œ,¡…∑*äKsˆÎÔ¸‚∂ø#√|Úˇ◊≈¸ˇåSΩªBl»Ã+I- KÃaXoZ|k]KÙÁu_5ß9¸<ˆuΩ#DIInC:kı°3ˆÊnOÂL&sãG%◊·dÅ
```

## .git/objects/79/86982687d75da66eb87f7a2588dd8d33f06898

> Fichier binaire non inclus (109 octets)

## .git/objects/7b/6628ade4bd1e5bf8024768ac26ff688ef83082

> Fichier binaire non inclus (407 octets)

## .git/objects/7c/d49f16753e2d5e87e82317cc62ba5ad25fcd1a

> Fichier binaire non inclus (1105 octets)

## .git/objects/7d/14ca0ceb65a100b65da50c13c4ed00423b7abd

> Fichier binaire non inclus (482 octets)

## .git/objects/7d/d15a1c616cdd4892d04ffee27e0b1f8e0240bb

> Fichier binaire non inclus (20962 octets)

## .git/objects/7e/8da60904b656207a5010b7027b47e651c85c0c

> Fichier binaire non inclus (110 octets)

## .git/objects/7f/caa3a05cd27baa703ea6e7989472bdc7f6e469

> Fichier binaire non inclus (1457 octets)

## .git/objects/81/692525cea495f555182927f8f55584d98020f9

> Fichier binaire non inclus (826 octets)

## .git/objects/83/fc2d1afb8f6583a9ac47503781a5f53b7e41c8

> Fichier binaire non inclus (54 octets)

## .git/objects/84/058f295333b603c0d827b8f70f6f199e64a3df

> Fichier binaire non inclus (76 octets)

## .git/objects/84/6efe57e7f1aba3621ccda2c424a11089a02b9b

> Fichier binaire non inclus (82 octets)

## .git/objects/84/dc0203f93f52853e23371eb3c46d49449ccf3e

```text
xïéMN!Ö]˜)j©qbÄAË6∆®—∏“dúE1C’˜ó+∏|/ﬂ˚°ZJ0ãπí∆3&KiÒ—%≥7—á5;euÙDC“¬ÙãçW
Jö¨bR*E´Sr^Ö@≥Û˚˚≈∫îˆó	79◊=”€À◊«Î;<F\{›NœßÇ˘rGµ<ÅˆŒπŸ,⁄√≠“JM√˜Ñˇúj£3wi(µ=@´aÎíSMT«n)(πÆπ√†Káo!ÀÁÆı∏§(∏n|ŸAM©≥‹L_`_Ó
```

## .git/objects/87/5968d8f81678a6e9a858eebe90a2b957759531

> Fichier binaire non inclus (143 octets)

## .git/objects/87/63c1d9aa8a5da1da10bdcac57f1fc5ebb09491

> Fichier binaire non inclus (55 octets)

## .git/objects/89/13571fad3aed041bd6f5a6fbb57446cf28ec84

```text
xïèIN≈0YÁΩGÑéÌx¯BàQÏ∏ÉáNbà„»±˜Çr˛'`[z%ΩÚ9•XÅqU†≤d‘IJ‰<êG70aGGN!G62°ßn∑Ö∂
ﬁXÕ‘Çq-ëî§s¡)°√‰x`ìDG¶≥≠.π¿˝À„˚€”+‹ªπÕs≤qÌ}N˜0()•¥—pçbw¶Á{ï˛-v!˚„ˆ#∑ã›œ±.Õ›˙º«5◊õ∏µ4_cﬁé>¯˘˙Üπ≈@`ÁK÷Ûﬂˆ‹
¨é6ø$[>ª_§D`é
```

## .git/objects/8a/f4c7ac97d6f232d2b6a4a1e6041d7ccb6ac8bb

> Fichier binaire non inclus (407 octets)

## .git/objects/8e/cc24ea2efe2596320b847a326edc2296cd2c7f

> Fichier binaire non inclus (4098 octets)

## .git/objects/91/f5acd81ec45ecd315e27d23991247bd92e1979

> Fichier binaire non inclus (119 octets)

## .git/objects/93/15641dac3ec2c85995ad959b95037cc83827c5

> Fichier binaire non inclus (1323 octets)

## .git/objects/96/47bb3ee733f230e94decddebc7276f600e2b32

> Fichier binaire non inclus (3238 octets)

## .git/objects/9d/c71497ae68a4e83a6d720c3ad2fe6abd38d86f

> Fichier binaire non inclus (305 octets)

## .git/objects/9d/cf46457a2239e54fbe4a9f00f2826996522e04

> Fichier binaire non inclus (168 octets)

## .git/objects/a2/91f3e344cf5c4aef3053ba2cb40e91b5fc39b3

> Fichier binaire non inclus (85 octets)

## .git/objects/a5/b9c675de4175687d7fff51d0dab61507fb118d

> Fichier binaire non inclus (48 octets)

## .git/objects/ab/ece7364f2ed1cdc9b439c4c5374ab3d9dfd837

```text
xÖêMN√0ÖY˚#Wb’Fi(K†@¨™
Ÿ…§Íÿ¡3Æ¯‚\«a¡éùı∆ÔÕ˚∆∫`a”‘ª—¥'s¿ΩÚf@∏›ì´òX¥:cd
~RÎj]’ZaGÚ´4u≥÷JÌ:—wË[Bﬁ´‹kI∂w–:B/ ôëJDS]È<øÂ7ﬂBL^h@%·D!/¯Ä?À¶UKË—Hä»y∫”}rNÔ·3<`$„Ë›LMcÏ¶ﬁˇtŸu∆Q,Oœ<ìÆlç·Ö+üÜï•yôŸ7∫H…‘'ü£—™O•ﬂ*eœ,f@3ºå±¸YW◊2Øå>ó=#åÅô¨√%Ü2úœ¿)¬|¶1î'Ù_…í#˘˛R?£@Ü¡
```

## .git/objects/ac/d9d78e43721d13834d5d0f2c10a60001a496b6

> Fichier binaire non inclus (215 octets)

## .git/objects/ad/7005ae32ba2e370bdae2663559142ddea29a5b

> Fichier binaire non inclus (15640088 octets)

## .git/objects/af/3573daae845bf3f07ef5299640f9c6f5af2201

> Fichier binaire non inclus (81 octets)

## .git/objects/b0/aaf9bab41421336a5b566acb37ad38096d48d3

> Fichier binaire non inclus (338 octets)

## .git/objects/b1/9ef482e0ef1f6fb118834ce1c56aec63821e1d

> Fichier binaire non inclus (83 octets)

## .git/objects/b5/3633d98bf09f65d7ffad41851e7dc00426fc7a

> Fichier binaire non inclus (84 octets)

## .git/objects/b7/6e4ffecf8364809c61229d2644b578f6a94157

> Fichier binaire non inclus (88 octets)

## .git/objects/b7/9c49648d6cba60ded7e23b9eccc64959b1031c

```text
xÖSQk€0ﬁkÙ+Cô]2∆÷≤¬É¶∞¡X
{(%ΩXr¨≈í=ÈƒJPGˇÿN±ú9ae~ë}wﬂ›}ü>ØÍfﬂûøj±ÿ‡Z¡Ô∆mîBõ∂qπeÖ€∂‘L}ÖÔ>úe(≈£ÒôA€V¡èÏªZkOI7¯%btı∫áC⁄Æ·ÓßoÏE¶ev'F_–ÆCú|ò¨SòK”åô4⁄^£ﬂZ«œÿ>Å”ú˝ë&W-ªqs¥≤VÓ?Ä™´ä˚-Õ¥•≥˜˚∆˚ó4¡«™%RD|Cá∆lo∫Ωn¨\âÖzxÏŸ∑ª™Ω^Bj0Ï¬«”ÿÓ∂µ^ãC5^ñ‡üuÊWV™˚Q^ˇò∞é…»ˆ´"îHQ°cí&%c·µ6äe2ÌÅî˝z}≤d:ÖKe;I1Û¬=?)®_k©,ÈR#Ø¨˛,Ê36%»êå@>?té=% `ãAß|‹ﬂ;ª≤j<Y4ÍñpÒ	?ôßH>£%≥◊'ãñmKeûù¯7'2;ÖøPF]*jµÃ«™–W±W˜´L¡r˘ÕÌjK*ècëSú=j|ª2zÃnˇ” Ç
```

## .git/objects/b7/e1fe52398f8f442711a0d0470a6318c0f10df3

> Fichier binaire non inclus (51 octets)

## .git/objects/be/0c220c16ac5a2391490b3f9618e5f9454345aa

> Fichier binaire non inclus (527 octets)

## .git/objects/bf/2fb9f446bf88683793ffeb730ffe33654b8be1

> Fichier binaire non inclus (4189 octets)

## .git/objects/c2/7fb11e9c9c3fe264ad1fdf894bbb8214c1d6a7

```text
x5è¡J1Ü=Á)~÷CÃÆ§‡IH+,“FÇ"R<§ŸX"1í-‘GÚ9|1ìÆ^Üôˇü˘ffÁiá˘Õ‚Ï^…5ˆv~Õµè.X∆^§zXˆ
ùéë1ıºAàüp!è⁄{‹ÌS4]¸#ˇU1—H‹ìlbÏÇ¢≥8∏0F{-í6e≈%ÚÒu2ª¸L;ıµuñs!ïZâß^n–ãæJ≈ê9˛|˜^ÿæ‡)gg2í±ò’7∫ŸÑ>Â.ˆÿñ#W¨óÿ6UnÆ–¸;Õ€/KT6
```

## .git/objects/c4/e58bd3837ab011a7decff3a984391061896366

```text
x}PAN√0‰øbÂSÇ⁄DH\®ƒA¢à«ŸƒVS€ÿkÑTıAºÉèa∑™‘Ø<ªûŸwìÌ‡Ú‚ÍÃ	π#Çpö1Ωq÷î¨‡©QDé≥FM*vµ¥õf‘f>Z£eæqV16D#°EäÓ›FB_Vpûfı“§ä∞eÖá≈5‰÷"NT&V·ÎáÂG…Öb"≈gêuJy‡ﬁZC¯MUf≤~jWØevS∑$(Ü’ÛlØ˜∏Âaﬂ‡‡b“_»w+“aE”¿}§ËA	”OË¡Ÿ&@(@è@^ò0`Jú¿†•“ËCvˆ∂j≥5mÁ˘Ìﬁ“ñ;ÍœàY√ˇ˛Ht§≠9Q<,Ó–Hµ~ΩtàqÈFf&ˆG˚>}®7‡Ÿé˝∆Râ)
```

## .git/objects/c4/f763aaa6fea4a56f4c36a8a0bd290b983f720b

> Fichier binaire non inclus (113 octets)

## .git/objects/c7/9719e52ab405c026c4d3c48899017ff6545063

> Fichier binaire non inclus (1043 octets)

## .git/objects/c8/7dd4961de5d2b0546460efbf71169d6fbde2c0

> Fichier binaire non inclus (475 octets)

## .git/objects/c9/a83108423860e7c4e872bdb748dfb3d2f60be9

> Fichier binaire non inclus (177 octets)

## .git/objects/ca/702652572c3cdc89b07f30a4f264c3cb78c40f

> Fichier binaire non inclus (2499 octets)

## .git/objects/cb/01ac9c40ec00fd41ff670bbc86735946ff34a9

```text
xïéKJ1Ego.H~ïO#¢¢8s‚¿ÒKÚR∫SëT
◊‰:‹òŸÇ£Áp‡∆Vk ≠∫ù¨BCë{ÓÃíÇîÜºŒ(î∞\$Œ≠6àh—≥/Ï¥p^®≈äåI!%ÆEH&/hrã’⁄ƒ,Eß„‹:Ï%æ<Ωø=ø¬}¬mo«˙∏V,◊ªÿÍkåq“s∑\pŒ&ù˜˝;dG9÷ﬂüki}Êòsâg\	Ì±è9ﬂ≠_¶£eK%‚†£√'Öè/4ÿXä
```

## .git/objects/d4/55b2fc342a290696ff081184fe4d5f18c87b7c

> Fichier binaire non inclus (541 octets)

## .git/objects/d4/ad699edca1759fdd98c08d197aea6fe8e7e202

> Fichier binaire non inclus (93 octets)

## .git/objects/d7/0c9ad928febd9287bcb9079f5857db1af09c06

> Fichier binaire non inclus (143 octets)

## .git/objects/d7/b6bdacc7db9681e4d5c83bbce813d0dae49244

> Fichier binaire non inclus (789 octets)

## .git/objects/dc/dd5be855c5c415d5179b9375b819fe9c3daf06

> Fichier binaire non inclus (925 octets)

## .git/objects/df/c5b6c588b2753fb511e78f198853374d94ddbc

> Fichier binaire non inclus (364 octets)

## .git/objects/e1/13dd3b75bdcf3d891223e4ee5a8e376a739d06

> Fichier binaire non inclus (68 octets)

## .git/objects/e2/0e09349e91e8cafe04f0cd25d2b449b3459170

```text
x=éKÉ0DªŒ),Ø[Ï{Ñ*B,)âQ Ø∑ØPW€o4c<hõˆñòxV®ñ0ø¨KyâJufuﬁˆ
¿∏X -Qéc∞˘81’)˝°¥Q*ñj‚Ï≤‰áØ¯ö∫Æe•cÙ´•+∏√ùﬁœ26äñÊ
.û πo'±Öõ˝Öó˘^Ç«˛∫Ae
```

## .git/objects/e2/b928d9112e1dc2bf98d9006e9a0bf2fc8cbf93

```text
xeRmk‘@ˆs~≈P¥ŸpπmäV§Òä»)“€oG°õ‹$¨∆›sw√!G˛ê√?÷ôÕΩYÁ√2;ÛÃ<ÛVu∂ÇwØﬂæ®≠ÒZ∑™aı⁄°H?∞·åüÈwüfe2¬VŒ˚’™%∫ˇ——9Ì¢˜(D’?Tãsl¥—A[CÅGi$„Ô~õZ§û^]áh∂›2)Eq5~˚<•X=∑dí#À$izSGV◊õœ∫≤"ÉM$çu :@J~^FÂ=iET'ìí—T+L√M≥∞oç ÚÃU@iÏZ––ÿÀ≤Á“ú¢ƒÛ¶`°œøLå_(by Ï¢ âvBüC2ËÑÜWT+	ÃfîNO)Ô«e3ñÖK¥“ÙZÒ∏∏æôzÄo‘˜Àu?¿î=ÄË<}ÓÒÁ _íMö!sÏ0;cŒ¡?µ«√æŸQí!ŸM›i4<Ékÿ/T~TŒit"mh-”÷^^≈≈yöè{Æ.)J´ŒK“i®_å«∫ß„Ã∂Á∏kË‰⁄.ënÍÔüêP!P$¬≤Ø[£:)Â	«ƒ:‰ZÈpoÔx]b3‰ –πfW€ª‡—FÀnfÙ±.ö∆`Áëoeº¶2≤Ú	Pû¸0
```

## .git/objects/e3/8fb63d9a6cac854ef98242004daeaa4ab8d694

> Fichier binaire non inclus (672 octets)

## .git/objects/e4/0d4a3fbc0649821696f81aeb021bb21fbda101

> Fichier binaire non inclus (81 octets)

## .git/objects/eb/3dda2764d74686f269b5551c818e3a4028218d

> Fichier binaire non inclus (229 octets)

## .git/objects/ec/d75440765f8c854fc0dd750e7f3249b9194b94

> Fichier binaire non inclus (830 octets)

## .git/objects/ee/d6e7bc041cbff29407b9c98dd75b168a5869fa

> Fichier binaire non inclus (44212 octets)

## .git/objects/ef/327c55831b06abfc2b987608b12db85728ffbf

```text
xUR›J√0ı∫Oq(»RÏJ+õHa^ÇW:‹@D$$]⁄÷¥$2∆^»◊≈L∂n≠ﬂEŒ9ﬂ¯¶·∏ÀÓØd›6⁄¢“mto+kî∫©av™†-øE«<’≠›˝ß®˜<ÛèLk)Ù¬ny¨D	ΩU¥îº!Q¿YŸhèA*h¶*A≤Y÷ë^¿bpÃê:¯Ë‡1cô∂‘WÂ%˛ —E‡£ >Ê$u6:Ãc0‹Ä_|='KÁ~çÃ˚aÊ≤É©ïÉêû ˆ¢≥µZ*K s˛±|~}˘¬õÎgÔö:`åΩ<@Z°ç{/E›ö˚A≈Ìõ…ìIy0a.?•ä’ÇRü>§¥fRQû≤Kª>.(ë àb´-÷L)±!#?‹q’‰”4ùf£Ã†„˙ ç[áõ‹`9§”Ù#<uŒwv›(¥˙˜«∆
ÃZ°¨¿j#+≈6Ií∏íœ√°ìw&Ì≤Y¯-ë„!QØË¿øy∂S
```

## .git/objects/f1/2d44b84ecbccc4b653c7fc833ad2dd7373db46

> Fichier binaire non inclus (1283 octets)

## .git/objects/f5/c1c09205f17763c8289c63c6d35fd4df64c6a8

> Fichier binaire non inclus (126 octets)

## .git/objects/f6/64bd3652c06f1111cd71c7133c9409dd882479

> Fichier binaire non inclus (54 octets)

## .git/objects/f8/786ea322dc83337b57a86a0e86612465bbd5b5

```text
xçê]KÎ@ÜΩÓØXFëVwõÏG@PhQA≠∆ê√°$ôù¥®Mw'Œ˘ıN*"Ωëﬁºs≥˚>ÛL˘÷îB'zo_\≠0‘∆ÆjªË˘trı0 ØO£á%e««_SèçˆE^îÂ≤ΩπC¡èH¨õ.äjQÑŒÚÒcYyqºh®ÕDÈW’‚ΩàØt(˛TÀ÷gBù®t(’D˝úﬂ‹ﬂı§˘ıÙ|2ÕOã˜∞fﬁÊ◊¸Õ‡ctäIÂ,T¡ë5u°ë ≠ìˇ]mç$SüIßFR€ëIï•⁄®Ø™Õ?∞ãŸ6®n2,pÅÀ!êNö8uË3–XIN2Ôa8¡≥#/zx‹&∆éè≥)1xƒÇhπúE”À&»+î,©1!ÿt˜Úx9ª›F≠ˇµãfïU≤_ô0ÈÜØÖç6∞πF«sWŒÌl2›¶¨FÙ
ﬂ"˝ÍSNd´_%>S/Õ‡
```

## .git/objects/fa/10cc7e0d046698a0c2b77544ed2a7167a49d25

> Fichier binaire non inclus (489 octets)

## .git/refs/heads/master

```text
5a8af5557d17aa3d37ee1e7841508587a0d49f19

```

## .git/refs/remotes/origin/master

```text
5a8af5557d17aa3d37ee1e7841508587a0d49f19

```

## .github/copilot-instructions.md

````markdown
# Copilot / Agent Instructions for Fibonnaci Benchmark

Quick, actionable notes to get an AI coding agent productive in this repository.

1) Big-picture architecture
- **Orchestrator (Go)**: `orchestrator/` runs an HTTP server + WebSocket hub and owns RabbitMQ orchestration. Entry: `orchestrator/cmd/server/main.go`.
- **Workers**: language-specific worker folders at `worker-go/`, `worker-node/`, `worker-python/`, `worker-rust/`. Each worker connects to RabbitMQ and participates in the benchmark.
- **Message bus**: RabbitMQ is the central integration point. Key code: `orchestrator/internal/orchestrator/rabbitmq.go`.
- **UI**: static files and client code under `orchestrator/internal/api/web/` and WebSocket hub in `orchestrator/internal/api/hub.go`.
- **Compose / orchestration**: `compose.yml` defines services, environment variables and network `fibo-benchmark-network`.

2) Core dataflows & conventions (must-follow)
- Worker registration: workers announce themselves to queue `isReady` (see `ListenForWorkers()` in `rabbitmq.go`).
- Per-worker result queues follow the pattern `results_{workerID}`. These are created as RabbitMQ streams (queue arg `x-queue-type: stream`). See the `QueueDeclare` call in `rabbitmq.go`.
- Admin sync: a durable fanout exchange `fibo_admin_exchange` is declared for synchronous broadcast across services.
- Durability: queues/exchanges are declared `durable: true` in orchestrator ‚Äî when adding or changing queues follow the same durability/stream semantics to avoid PRECONDITION_FAILED errors.
- Concurrency: orchestrator uses `Engine.Mu` to guard `Channel`/`Workers`. Follow locking patterns when mutating shared Engine state.

3) Important files to reference
- RabbitMQ orchestration: `orchestrator/internal/orchestrator/rabbitmq.go`
- Engine and types: `orchestrator/internal/orchestrator/types.go`
- WebSocket hub & routes: `orchestrator/internal/api/hub.go`, `orchestrator/internal/api/routes.go` and `orchestrator/internal/api/web/`
- Docker compose: `compose.yml` (service names, env var names like `AMQP_URL_*`, `REDIS_URL`)
- Worker examples: `worker-go/`, `worker-node/`, `worker-python/`, `worker-rust/` (follow folder layout and Dockerfiles)
- Proto definitions: `proto/sync.proto` (check if type generation is needed before cross-language changes)

4) Developer workflows (how to build/run/debug)
- Full stack (recommended): from repo root
  - Build & start: `docker compose up -d --force-recreate --build --remove-orphans`
  - Tail orchestrator logs: `docker compose logs -f orchestrator`
  - Start a single service: `docker compose up -d orchestrator` or `docker compose up -d fibo-go`
- Local Go iteration (without Docker): set `AMQP_URL` locally then run:
  - `go run ./orchestrator/cmd/server` (or `go build ./...` then run binary)
  - The orchestrator listens on `:8080` and exposes health at `/health`.
- Worker iteration:
  - Go worker: build/run under `worker-go/` (`go run ./cmd/worker`)
  - Node worker: `node worker-node/index.js`
  - Python worker: `python3 worker-python/main.py`
  - Rust worker: `cargo build` / `cargo run` in `worker-rust/`

5) Environment variables and defaults
- `AMQP_URL` (or variants in `compose.yml`: `AMQP_URL_LEADER`, `AMQP_URL_GO`, `AMQP_URL_NODE`, `AMQP_URL_PYTHON`, `AMQP_URL_RUST`) ‚Äî orchestrator falls back to an embedded URL in `orchestrator/cmd/server/main.go` if unset.
- `REDIS_URL` referenced in compose for caching/coordination.

6) When changing message schemas or queues
- Update `proto/sync.proto` first if the change affects message formats used by multiple languages; regenerate stubs if the repo uses codegen (no generator found in repo‚Äîconfirm with maintainers).
- Keep queue names and durability consistent with `rabbitmq.go` (`isReady`, `results_{id}`, `fibo_admin_exchange`). Streams require declaring `x-queue-type: stream`.

7) Code patterns the agent should follow
- Keep resilient reconnect logic intact: `Engine.InitRabbitMQ` ‚Üí `handleReconnect` is responsible for reconnection and re-declaring infra.
- Use `Engine.Mu` when accessing `Engine.Channel` or `Engine.Workers`.
- When adding a worker implementation, mirror the `results_{id}` creation and support the `isReady` registration flow.
- Prefer the existing logging style (`log.Printf` with tags like `[RMQ]`, `[SYNC]`, `[WORKER]`).

8) Quick examples (copyable)
- Declare a stream results queue (Go):
```go
args := amqp.Table{"x-queue-type":"stream"}
_, _ = ch.QueueDeclare("results_<id>", true, false, false, false, args)
```
- Broadcast to UI via hub:
```go
hub.BroadcastMessage(map[string]interface{"type": "WORKER_JOIN", "data": reg})
```

9) What the agent should ask the maintainers (if unclear)
- Are stubs generated from `proto/sync.proto` in CI or manually? Where is generation configured?
- Which worker languages are currently considered canonical for benchmarks (the compose file comments indicate some are commented out)?
- Any cluster-specific assumptions (k3s, network addresses) that must be preserved when tweaking reconnect logic?

If anything above is unclear or you want more detail in a specific area (proto/codegen, worker onboarding, Docker/CI), tell me which section to expand. Thanks!

````

## .gitignore

```text
# Binaires Go
go/pb/*.go
fibo
fibo-go

# Node
node_modules/
npm-debug.log

# Python
__pycache__/
*.py[cod]

# OS
.DS_Store

```

## compose.yml

```yaml
services:
  redis:
    image: redis:7.4-alpine
    container_name: fibo-redis
    ports:
      - "6379:6379"
    command: redis-server --save 60 1 --loglevel warning
    networks:
      - fibo-net

  orchestrator:
    build:
      context: ./orchestrator
      dockerfile: Dockerfile
    container_name: fibo-orchestrator
    ports:
      - "8080:8080"
    environment:
      - AMQP_URL=${AMQP_URL_LEADER}
      - REDIS_URL=${REDIS_URL}
      - GIN_MODE=release
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 5s
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 2G
    networks:
      - fibo-net

  fibo-go:
    build:
      context: ./worker-go
      dockerfile: Dockerfile
    container_name: worker-go
    environment:
      - AMQP_URL=${AMQP_URL_GO}
      - REDIS_URL=${REDIS_URL}
    # depends_on:
    #   orchestrator:
    #     condition: service_healthy
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks:
      - fibo-net

  # fibo-rust:
  #   build:
  #     context: ./worker-rust
  #     dockerfile: Dockerfile
  #   container_name: worker-rust
  #   environment:
  #     - AMQP_URL=${AMQP_URL_RUST}
  #   depends_on:
  #     orchestrator:
  #       condition: service_healthy
  #   deploy:
  #     resources:
  #       limits:
  #         cpus: '5.0'
  #         memory: 5G
  #   networks:
  #     - fibo-net

  # fibo-python:
  #   build:
  #     context: ./worker-python
  #     dockerfile: Dockerfile
  #   container_name: worker-python
  #   environment:
  #     - AMQP_URL=${AMQP_URL_PYTHON}
  #   depends_on:
  #     orchestrator:
  #       condition: service_healthy
  #   deploy:
  #     resources:
  #       limits:
  #         cpus: '5.0'
  #         memory: 5G
  #   networks:
  #     - fibo-net

  # fibo-node:
  #   build:
  #     context: ./worker-node
  #     dockerfile: Dockerfile
  #   container_name: worker-node
  #   environment:
  #     - AMQP_URL=${AMQP_URL_NODE}
  #   depends_on:
  #     orchestrator:
  #       condition: service_healthy
  #   deploy:
  #     resources:
  #       limits:
  #         cpus: '5.0'
  #         memory: 5G
  #   networks:
  #     - fibo-net

networks:
  fibo-net:
    driver: bridge
    name: fibo-benchmark-network
```

## orchestrator/.air.toml

```toml
root = "."
tmp_dir = "tmp"

[build]
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server/main.go"
  delay = 1000
  exclude_dir = ["web", "tmp", "vendor"]
  include_ext = ["go", "tpl", "tmpl", "html"]
```

## orchestrator/Dockerfile

```text
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache ca-certificates git
WORKDIR /src

# Cache des d√©pendances
COPY go.mod go.su* ./
RUN go mod download

# Copie de tout le code pour permettre l'embed
COPY . .

# Build statique (Correction du chemin vers cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/orchestrator ./cmd/server/main.go

# --- Stage 2: Image Finale ---
FROM scratch
# Import des certificats pour RabbitMQ TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /

# Le binaire contient l'IHM gr√¢ce √† go:embed
COPY --from=builder /app/orchestrator .

EXPOSE 8080
ENTRYPOINT ["./orchestrator"]
```

## orchestrator/cmd/server/main.go

```go
package main

import (
	"context"
	"fibo-orchestrateur/internal/api"
	"fibo-orchestrateur/internal/orchestrator"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. Initialisation du Hub WebSocket
	hub := api.NewHub()
	go hub.Run()

	// 2. Configuration de l'URL RabbitMQ via .env
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://bench_leader:qsd65f4c98dc7fd9s87ga6fsd5g4zsdrf9g879dfs7g@192.168.1.12:5672/benchmarks"
	}

	// 3. Lancement de l'Engine (RabbitMQ + Orchestration)
	// NewEngine lance d√©j√† la boucle de reconnexion en interne
	orch := orchestrator.NewEngine(amqpURL, hub)

	// 4. Configuration du serveur HTTP avec Graceful Shutdown
	router := api.SetupRouter(orch, hub)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Lancement du serveur dans une goroutine pour ne pas bloquer le thread principal
	go func() {
		log.Println("[API] Serveur d√©marr√© sur :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[CRITICAL] Erreur serveur Web: %v", err)
		}
	}()

	// 5. Gestion des Signaux d'Arr√™t (SIGINT, SIGTERM)
	// Indispensable pour k3s lors d'un red√©ploiement
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Blocage jusqu'√† r√©ception d'un signal
	sig := <-quit
	log.Printf("[MAIN] Signal %v re√ßu. D√©but de la proc√©dure d'arr√™t propre...", sig)

	// 6. Proc√©dure de fermeture Graceful (Timeout de 5 secondes)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Arr√™t du serveur HTTP (ne prend plus de nouvelles requ√™tes)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[ERROR] Arr√™t forc√© du serveur HTTP:", err)
	}

	// Fermeture des ressources critiques
	log.Println("[RMQ] Fermeture des connexions RabbitMQ...")
	if orch.Conn != nil {
		orch.Conn.Close()
	}

	log.Println("[MAIN] Orchestrateur arr√™t√© proprement. Bye!")
}

```

## orchestrator/go.mod

```text
module fibo-orchestrateur

go 1.25.1

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/gorilla/websocket v1.5.3
	github.com/streadway/amqp v1.1.0
)

require (
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

```

## orchestrator/go.sum

```text
github.com/bytedance/sonic v1.14.0 h1:/OfKt8HFw0kh2rj8N0F6C/qPGRESq0BbaNZgcNXXzQQ=
github.com/bytedance/sonic v1.14.0/go.mod h1:WoEbx8WTcFJfzCe0hbmyTGrfjt8PzNEBdxlNUO24NhA=
github.com/bytedance/sonic/loader v0.3.0 h1:dskwH8edlzNMctoruo8FPTJDF3vLtDT0sXZwvZJyqeA=
github.com/bytedance/sonic/loader v0.3.0/go.mod h1:N8A3vUdtUebEY2/VQC0MyhYeKUFosQU6FxH2JmUe6VI=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/gabriel-vasile/mimetype v1.4.8 h1:FfZ3gj38NjllZIeJAmMhr+qKL8Wu+nOoI3GqacKw1NM=
github.com/gabriel-vasile/mimetype v1.4.8/go.mod h1:ByKUIKGjh1ODkGM1asKUbQZOLGrPjydw3hYPU2YU9t8=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.11.0 h1:OW/6PLjyusp2PPXtyxKHU0RbX6I/l28FTdDlae5ueWk=
github.com/gin-gonic/gin v1.11.0/go.mod h1:+iq/FyxlGzII0KHiBGjuNn4UNENUlKbGlNmc+W50Dls=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.27.0 h1:w8+XrWVMhGkxOaaowyKH35gFydVHOvC0/uWoy2Fzwn4=
github.com/go-playground/validator/v10 v10.27.0/go.mod h1:I5QpIEbmr8On7W0TktmJAumgzX4CA1XNl4ZmDuVHKKo=
github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
github.com/goccy/go-yaml v1.18.0 h1:8W7wMFS12Pcas7KU+VVkaiCng+kG8QiFeFwzFb+rwuw=
github.com/goccy/go-yaml v1.18.0/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/gorilla/websocket v1.5.3 h1:saDtZ6Pbx/0u+bgYQ3q96pZgCzfhKXGPqt7kZ72aNNg=
github.com/gorilla/websocket v1.5.3/go.mod h1:YR8l580nyteQvAITg2hZ9XVh4b55+EU/adAjf1fMHhE=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 h1:ZqeYNhU3OHLH3mGKHDcjJRFFRrJa6eAM5H+CtDdOsPc=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/quic-go/qpack v0.5.1 h1:giqksBPnT/HDtZ6VhtFKgoLOWmlyo9Ei6u9PqzIMbhI=
github.com/quic-go/qpack v0.5.1/go.mod h1:+PC4XFrEskIVkcLzpEkbLqq1uCoxPhQuvK5rH1ZgaEg=
github.com/quic-go/quic-go v0.54.0 h1:6s1YB9QotYI6Ospeiguknbp2Znb/jZYjZLRXn9kMQBg=
github.com/quic-go/quic-go v0.54.0/go.mod h1:e68ZEaCdyviluZmy44P6Iey98v/Wfz6HCjQEm+l8zTY=
github.com/streadway/amqp v1.1.0 h1:py12iX8XSyI7aN/3dUT8DFIDJazNJsVJdxNVEpnQTZM=
github.com/streadway/amqp v1.1.0/go.mod h1:WYSrTEYHOXHd0nwFeUXAe2G2hRnQT+deZJJf88uS9Bg=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/ugorji/go/codec v1.3.0 h1:Qd2W2sQawAfG8XSvzwhBeoGq71zXOC/Q1E9y/wUcsUA=
github.com/ugorji/go/codec v1.3.0/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
go.uber.org/mock v0.5.0 h1:KAMbZvZPyBPWgD14IrIQ38QCyjwpvVVV6K/bHl1IwQU=
go.uber.org/mock v0.5.0/go.mod h1:ge71pBPLYDk7QIi1LupWxdAykm7KIEFchiOqd6z7qMM=
golang.org/x/arch v0.20.0 h1:dx1zTU0MAE98U+TQ8BLl7XsJbgze2WnNKF/8tGp/Q6c=
golang.org/x/arch v0.20.0/go.mod h1:bdwinDaKcfZUGpH09BB7ZmOfhalA8lQdzl62l8gGWsk=
golang.org/x/crypto v0.40.0 h1:r4x+VvoG5Fm+eJcxMaY8CQM7Lb0l1lsmjGBQ6s8BfKM=
golang.org/x/crypto v0.40.0/go.mod h1:Qr1vMER5WyS2dfPHAlsOj01wgLbsyWtFn/aY+5+ZdxY=
golang.org/x/mod v0.25.0 h1:n7a+ZbQKQA/Ysbyb0/6IbB1H/X41mKgbhfv7AfG/44w=
golang.org/x/mod v0.25.0/go.mod h1:IXM97Txy2VM4PJ3gI61r1YEk/gAj6zAHN3AdZt6S9Ww=
golang.org/x/net v0.42.0 h1:jzkYrhi3YQWD6MLBJcsklgQsoAcw89EcZbJw8Z614hs=
golang.org/x/net v0.42.0/go.mod h1:FF1RA5d3u7nAYA4z2TkclSCKh68eSXtiFwcWQpPXdt8=
golang.org/x/sync v0.16.0 h1:ycBJEhp9p4vXvUZNszeOq0kGTPghopOL8q0fq3vstxw=
golang.org/x/sync v0.16.0/go.mod h1:1dzgHSNfp02xaA81J2MS99Qcpr2w7fw1gpm99rleRqA=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.35.0 h1:vz1N37gP5bs89s7He8XuIYXpyY0+QlsKmzipCbUtyxI=
golang.org/x/sys v0.35.0/go.mod h1:BJP2sWEmIv4KK5OTEluFJCKSidICx8ciO85XgH3Ak8k=
golang.org/x/text v0.27.0 h1:4fGWRpyh641NLlecmyl4LOe6yDdfaYNrGb2zdfo4JV4=
golang.org/x/text v0.27.0/go.mod h1:1D28KMCvyooCX9hBiosv5Tz/+YLxj0j7XhWjpSUF7CU=
golang.org/x/tools v0.34.0 h1:qIpSLOxeCYGg9TrcJokLBG4KFA6d795g0xkBkiESGlo=
golang.org/x/tools v0.34.0/go.mod h1:pAP9OwEaY1CAW3HOmg3hLZC5Z0CCmzjAF2UQMSqNARg=
google.golang.org/protobuf v1.36.9 h1:w2gp2mA27hUeUzj9Ex9FBjsBm40zfaDtEWow293U7Iw=
google.golang.org/protobuf v1.36.9/go.mod h1:fuxRtAxBytpl4zzqUh6/eyUujkJdNiuEkXntxiD/uRU=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=

```

## orchestrator/internal/api/hub.go

```go
package api

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan interface{}
	Mu        sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
	}
}

func (h *Hub) Run() {
	for msg := range h.Broadcast {
		h.Mu.Lock()
		for client := range h.Clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(h.Clients, client)
			}
		}
		h.Mu.Unlock()
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	hub.Mu.Lock()
	hub.Clients[conn] = true
	hub.Mu.Unlock()
}

func (h *Hub) BroadcastMessage(msg interface{}) {
	h.Broadcast <- msg
}

```

## orchestrator/internal/api/routes.go

```go
package api

import (
	"embed"
	"fibo-orchestrateur/internal/orchestrator"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed all:web
var webAssets embed.FS

func SetupRouter(orch *orchestrator.Engine, hub *Hub) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Endpoint utilis√© par le healthcheck du Compose
	r.GET("/health", func(c *gin.Context) {
		orch.Mu.Lock()
		// V√©rification de l'√©tat du canal RabbitMQ
		isRMQConnected := orch.Channel != nil && !orch.Conn.IsClosed()
		orch.Mu.Unlock()

		if isRMQConnected {
			c.JSON(http.StatusOK, gin.H{"status": "UP"})
		} else {
			// Retourne 503 pour que curl -f √©choue
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		}
	})

	webRoot, _ := fs.Sub(webAssets, "web")
	r.GET("/", func(c *gin.Context) {
		index, _ := fs.ReadFile(webRoot, "index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	r.StaticFS("/static", http.FS(webRoot))
	r.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)

		// On envoie imm√©diatement la liste des workers d√©j√† connus
		orch.Mu.Lock()
		for _, w := range orch.Workers {
			hub.BroadcastMessage(map[string]interface{}{
				"type": "WORKER_JOIN",
				"data": w,
			})
		}
		orch.Mu.Unlock()
	})
	r.POST("/run", func(c *gin.Context) {
		var req struct {
			Handler string                 `json:"handler"`
			Params  map[string]interface{} `json:"params"`
		}
		if err := c.BindJSON(&req); err == nil {
			orch.StartTask(req.Handler, req.Params)
			c.Status(http.StatusAccepted)
		}
	})
	return r
}

```

## orchestrator/internal/api/web/index.html

```html
<!DOCTYPE html>
<html lang="fr" class="bg-slate-900 text-slate-100">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Benchmark IA Orchestrator</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body class="min-h-screen">
    <div id="app" class="p-6 lg:p-12">
        <header class="flex justify-between items-center mb-8 border-b border-slate-800 pb-6">
            <h1 class="text-2xl font-black tracking-tighter text-blue-500 italic">BENCHMARK_LEADER_V1</h1>
            <div id="ws-status" class="w-4 h-4 rounded-full bg-red-500 shadow-[0_0_10px_red] transition-all duration-500"></div>
        </header>

        <div class="grid grid-cols-12 gap-8">
            <aside class="col-span-12 lg:col-span-3 space-y-6">
                <section class="bg-slate-800 p-6 rounded-2xl border border-slate-700 shadow-xl">
                    <h2 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Commandes</h2>
                    <button id="btn-fibo" class="w-full bg-blue-600 hover:bg-blue-500 py-4 rounded-xl font-black text-sm tracking-widest transition-all active:scale-95 shadow-lg shadow-blue-900/20">
                        LAUNCH_FIBONACCI
                    </button>
                </section>

                <section class="bg-slate-800 p-6 rounded-2xl border border-slate-700 shadow-xl">
                    <h2 class="text-xs font-bold uppercase tracking-widest text-slate-500 mb-4">Cluster Workers</h2>
                    <div id="worker-list" class="space-y-3">
                        </div>
                </section>
            </aside>

            <main class="col-span-12 lg:col-span-9 bg-slate-800 p-8 rounded-[2.5rem] border border-slate-700 shadow-2xl overflow-hidden relative">
                <div class="absolute top-0 right-0 p-4 text-[10px] font-mono text-slate-600 uppercase">Realtime_Stream_Active</div>
                <canvas id="mainChart" class="w-full"></canvas>
            </main>
        </div>
    </div>

    <script type="module" src="/static/js/app.js"></script>
</body>
</html>
```

## orchestrator/internal/api/web/js/app.js

```javascript
import { updateWorkerList, updateStatus ,updateHealthData} from './ui.js';
import { initChart, addDataToChart } from './charts.js';

const socket = new WebSocket(`ws://${window.location.host}/ws`);

socket.onopen = () => updateStatus(true);
socket.onclose = () => updateStatus(false);

socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    
    switch(msg.type) {
        case "WORKER_JOIN":
            updateWorkerList(msg.data);
            break;
        case "RESULT":
            addDataToChart(msg.data);
            updateWorkerMetrics(msg.data.worker_id, msg.data.metadata); // Mise √† jour IHM temps r√©el
            break;
        case "HEALTH_UPDATE":
            updateHealthData(msg.data);
            break;
    }
};

document.getElementById('btn-fibo').onclick = async () => {
    await fetch('/run', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            handler: "fibonacci",
            params: { series: 5, limit: 400000 }
        })
    });
};

initChart();
```

## orchestrator/internal/api/web/js/charts.js

```javascript
let chart;
const workerColors = {
    'rust': '#f97316',
    'go': '#00add8',
    'node': '#84cc16',
    'python': '#3b82f6'
};

export function initChart() {
    const ctx = document.getElementById('mainChart').getContext('2d');
    chart = new Chart(ctx, {
        type: 'line',
        data: { datasets: [] },
        options: {
            animation: false,
            scales: {
                x: { type: 'linear', grid: { color: '#1e293b' } },
                y: { grid: { color: '#1e293b' } }
            },
            plugins: { legend: { position: 'bottom' } }
        }
    });
}

export function addDataToChart(res) {
    let dataset = chart.data.datasets.find(d => d.id === res.worker_id);
    
    if (!dataset) {
        dataset = {
            id: res.worker_id,
            label: `${res.handler} - ${res.worker_id.substring(0,6)}`,
            data: [],
            borderColor: workerColors[res.handler] || '#fff',
            borderWidth: 2,
            pointRadius: 0
        };
        chart.data.datasets.push(dataset);
    }

    dataset.data.push({ x: res.index, y: res.timestamp % 10000 }); // Exemple de m√©trique
    if (dataset.data.length > 100) dataset.data.shift();
    chart.update('none');
}
```

## orchestrator/internal/api/web/js/ui.js

```javascript
export function updateStatus(connected) {
    const indicator = document.getElementById('ws-status');
    if (connected) {
        indicator.classList.replace('bg-red-500', 'bg-emerald-500');
        indicator.classList.replace('shadow-[0_0_10px_red]', 'shadow-[0_0_10px_#10b981]');
    } else {
        indicator.classList.replace('bg-emerald-500', 'bg-red-500');
        indicator.classList.replace('shadow-[0_0_10px_#10b981]', 'shadow-[0_0_10px_red]');
    }
}

export function updateWorkerMetrics(workerId, metadata) {
    const el = document.getElementById(`worker-${workerId}`);
    if (!el) return;

    const statsContainer = el.querySelector('.worker-stats') || document.createElement('div');
    if (!el.querySelector('.worker-stats')) {
        statsContainer.className = "worker-stats mt-3 pt-2 border-t border-slate-600/30 grid grid-cols-2 gap-2";
        el.appendChild(statsContainer);
    }

    statsContainer.innerHTML = `
        <div class="flex flex-col">
            <span class="text-[9px] text-slate-500 uppercase">Load</span>
            <span class="text-xs font-mono text-blue-400">${metadata.cpu} gor</span>
        </div>
        <div class="flex flex-col text-right">
            <span class="text-[9px] text-slate-500 uppercase">Mem</span>
            <span class="text-xs font-mono text-emerald-400">${metadata.ram}MB</span>
        </div>
        <div class="col-span-2 text-center text-[10px] text-slate-400 font-mono mt-1">
            Data Stream: ${(metadata.net / 1024).toFixed(2)} KB
        </div>
    `;
}

export function updateWorkerList(worker) {
    const list = document.getElementById('worker-list');
    const id = `worker-${worker.id}`;
    if (document.getElementById(id)) return;

    const el = document.createElement('div');
    el.id = id;
    // Ajout des classes pour l'expansion et l'interactivit√© 
    el.className = "p-4 bg-slate-800/80 rounded-2xl border border-slate-700 cursor-pointer transition-all duration-300 hover:border-blue-500 overflow-hidden mb-4";
    
    el.onclick = () => {
        // Toggle pour agrandir la tuile sur tout le conteneur si besoin
        el.classList.toggle('ring-2');
        el.classList.toggle('ring-blue-500');
        el.querySelector('.metrics-grid').classList.toggle('hidden');
    };

    el.innerHTML = `
        <div class="flex justify-between items-center">
            <div class="flex flex-col">
                <span class="text-xs font-black text-blue-400 italic tracking-tighter">${worker.language.toUpperCase()}</span>
                <span class="text-[9px] font-mono text-slate-500">${worker.id.substring(0,16)}</span>
            </div>
            <div class="status-dot w-2 h-2 rounded-full bg-emerald-500 shadow-[0_0_8px_emerald]"></div>
        </div>
        
        <div class="metrics-grid hidden mt-4 pt-4 border-t border-slate-700 grid grid-cols-2 gap-4">
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">RAM Usage</p>
                <p class="ram-val text-sm font-black text-white">-- MB</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">CPU Threads</p>
                <p class="cpu-val text-sm font-black text-white">-- gor</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">Net I/O</p>
                <p class="net-val text-xs font-mono text-slate-300">--</p>
            </div>
            <div class="metric-box bg-slate-900/50 p-3 rounded-xl">
                <p class="text-[8px] text-slate-500 uppercase font-bold">Disk Activity</p>
                <p class="disk-val text-xs font-mono text-slate-300">--</p>
            </div>
        </div>
    `;
    list.appendChild(el);
}

export function updateHealthData(data) {
    const el = document.getElementById(`worker-${data.worker_id}`);
    if (!el) return;
    console.log('Updating health data for', data.worker_id, data);
    // Mise √† jour des valeurs en temps r√©el 
    el.querySelector('.ram-val').innerText = `${data.ram} MB`;
    el.querySelector('.cpu-val').innerText = `${data.cpu} thr`;
    el.querySelector('.net-val').innerText = data.net_io;
    el.querySelector('.disk-val').innerText = data.disk_io;
}
```

## orchestrator/internal/orchestrator/engine.go

```go
package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Engine struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Workers map[string]WorkerRegistration
	Mu      sync.Mutex
	Hub     UIHub
}

type UIHub interface {
	BroadcastMessage(msg interface{})
}

func (o *Engine) safeGo(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[CRITICAL] Panic captur√© dans l'orchestrateur: %v", r)
				// Ici, on pourrait envoyer une alerte √† ton Sentry ou Loki
			}
		}()
		fn()
	}()
}

func NewEngine(amqpURL string, hub UIHub) *Engine {
	e := &Engine{
		Workers: make(map[string]WorkerRegistration),
		Hub:     hub,
	}
	e.InitRabbitMQ(amqpURL)
	return e
}

// StartTask diffuse l'ordre de calcul avec synchronisation temporelle (Phase 3)
func (o *Engine) StartTask(handler string, params map[string]interface{}) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if o.Channel == nil {
		log.Println("[WARN] Abandon : Orchestrateur non connect√© √† RabbitMQ")
		return
	}

	task := AdminTask{
		TaskID:  fmt.Sprintf("T-%d", time.Now().Unix()),
		Handler: handler,
		StartAt: time.Now().Add(5 * time.Second).Unix(), // Barri√®re T+5s pour tous les workers
		Params:  params,
	}

	body, _ := json.Marshal(task)
	err := o.Channel.Publish(
		"fibo_admin_exchange",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("[RMQ] √âchec envoi Task: %v", err)
	}
}

// ConsumeWorkerResults traite les messages entrants de chaque langage

func (o *Engine) ConsumeWorkerResults(queueName string) {
	// 1. On cr√©e un canal D√âDI√â pour ce flux Stream sp√©cifique
	ch, err := o.Conn.Channel()
	if err != nil {
		log.Printf("[ERROR] Impossible de cr√©er un canal pour %s: %v", queueName, err)
		return
	}
	defer ch.Close() // Fermeture propre si le worker dispara√Æt

	// 2. Configuration obligatoire pour Stream (Prefetch > 0)
	if err := ch.Qos(100, 0, false); err != nil {
		log.Printf("[RMQ] Erreur QoS pour %s: %v", queueName, err)
		return
	}

	// 3. Consommation avec offset 'next'
	msgs, err := ch.Consume(
		queueName, "",
		false, // auto-ack: false obligatoire
		false, false, false,
		amqp.Table{"x-stream-offset": "next"},
	)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume %s: %v", queueName, err)
		return
	}

	for d := range msgs {
		var res WorkerResult
		if err := json.Unmarshal(d.Body, &res); err == nil {
			o.BroadcastToUI("RESULT", res)
			d.Ack(false) // Ack manuel sur canal d√©di√©
		}
	}
}

func (o *Engine) BroadcastToUI(msgType string, data interface{}) {
	if o.Hub != nil {
		o.Hub.BroadcastMessage(map[string]interface{}{
			"type": msgType,
			"data": data,
		})
	}
}

func (o *Engine) StartGarbageCollector() {
	ticker := time.NewTicker(30 * time.Second)
	o.safeGo(func() {
		for range ticker.C {
			o.Mu.Lock()
			now := time.Now().Unix()
			for id, worker := range o.Workers {
				if now-worker.LastSeen > 60 {
					log.Printf("[GC] Worker %s inactif. Nettoyage...", id)

					// 1. Suppression de la queue durable sur RabbitMQ
					if o.Channel != nil {
						queueName := "results_" + id
						_, err := o.Channel.QueueDelete(queueName, false, false, false)
						if err != nil {
							log.Printf("[GC] Erreur suppression queue %s: %v", queueName, err)
						}
					}

					// 2. Nettoyage m√©moire et UI
					delete(o.Workers, id)
					o.BroadcastToUI("WORKER_LEAVE", map[string]string{"id": id})
				}
			}
			o.Mu.Unlock()
		}
	})
}

```

## orchestrator/internal/orchestrator/rabbitmq.go

```go
package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// InitRabbitMQ amorce la boucle de connexion r√©siliente
func (o *Engine) InitRabbitMQ(url string) {
	go o.handleReconnect(url)
}

// handleReconnect assure la survie de la connexion sur ton cluster k3s
func (o *Engine) handleReconnect(url string) {
	for {
		log.Printf("[RMQ] Tentative de connexion √† %s...", url)
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Printf("[RMQ] √âchec Dial: %v. Re-tentative dans 10s...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		o.Conn = conn
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[RMQ] √âchec Canal: %v", err)
			conn.Close()
			time.Sleep(10 * time.Second)
			continue
		}

		o.Mu.Lock()
		o.Channel = ch
		o.Mu.Unlock()

		if err := o.setupInfrastructure(); err != nil {
			log.Printf("[RMQ] Erreur Infrastructure: %v", err)
			conn.Close()
			time.Sleep(10 * time.Second)
			continue
		}

		log.Println("[RMQ] Connect√© et infrastructure pr√™te.")

		// --- Surveillance double : Connexion + Canal ---
		notifyConnClose := make(chan *amqp.Error)
		o.Conn.NotifyClose(notifyConnClose)

		notifyChanClose := make(chan *amqp.Error)
		o.Channel.NotifyClose(notifyChanClose)

		// Lancement des √©coutes
		go o.ListenForWorkers()
		go o.ListenForHealth()

		// Blocage intelligent
		select {
		case err := <-notifyConnClose:
			log.Printf("[RMQ] Connexion perdue (vhost/server): %v", err)
		case err := <-notifyChanClose:
			log.Printf("[RMQ] Canal perdu (queue/protocol error): %v", err)
		}

		// Nettoyage avant reconnexion
		o.Mu.Lock()
		o.Channel = nil
		o.Mu.Unlock()
		o.Conn.Close()

		log.Println("[RMQ] Nettoyage effectu√©. Re-tentative de connexion...")
		time.Sleep(5 * time.Second)
	}
}

func (o *Engine) setupInfrastructure() error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	// Exchange Fanout pour la synchronisation synchrone (Phase 3)
	err := o.Channel.ExchangeDeclare(
		"fibo_admin_exchange",
		"fanout",
		true, // Durable pour la persistance dans k3s
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Queue isReady pour Phase 1
	// Changement : durable: true pour √©viter l'erreur PRECONDITION_FAILED (406)
	_, err = o.Channel.QueueDeclare(
		"isReady",
		true, // Durable : doit correspondre √† la queue existante sur ton RMQ
		false,
		false,
		false,
		nil,
	)
	// D√©claration de la queue de sant√© (Classic Durable)
	_, err = o.Channel.QueueDeclare("worker_health", true, false, false, false, nil)
	return err
}

func (o *Engine) ListenForHealth() {
	o.Mu.Lock()
	ch := o.Channel
	o.Mu.Unlock()

	// S√©curit√© suppl√©mentaire : Nil check
	if ch == nil {
		log.Println("[ERROR] ListenForHealth abandonn√© : canal nil")
		return
	}

	msgs, err := ch.Consume("worker_health", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume health: %v", err)
		return
	}

	for d := range msgs {
		var healthData struct {
			WorkerID string `json:"worker_id"`
			RAM      int    `json:"ram"`
			CPU      int    `json:"cpu"`
		}
		if err := json.Unmarshal(d.Body, &healthData); err == nil {
			o.Mu.Lock()
			if w, ok := o.Workers[healthData.WorkerID]; ok {
				w.LastSeen = time.Now().Unix() // Rafra√Æchissement
				o.Workers[healthData.WorkerID] = w
			}
			o.Mu.Unlock()
			o.BroadcastToUI("HEALTH_UPDATE", healthData)
		}
	}
}

func (o *Engine) ListenForWorkers() {
	o.Mu.Lock()
	if o.Channel == nil {
		o.Mu.Unlock()
		return
	}
	ch := o.Channel // On extrait le canal pour ne pas garder le lock durant le Consume
	o.Mu.Unlock()

	msgs, err := ch.Consume("isReady", "", true, false, false, false, nil)
	if err != nil {
		log.Printf("[RMQ] Erreur Consume: %v", err)
		return
	}

	for d := range msgs {
		var reg WorkerRegistration
		if err := json.Unmarshal(d.Body, &reg); err != nil {
			continue
		}

		// --- SCOPE DE VERROUILLAGE ATOMIQUE ---
		o.Mu.Lock()
		reg.LastSeen = time.Now().Unix() // Initialisation du timestamp pour le GC
		o.Workers[reg.ID] = reg
		o.Mu.Unlock() // ON LIB√àRE IMM√âDIATEMENT

		// Cr√©ation de la queue Stream sur un canal d√©di√© pour √©viter l'erreur 503
		resQueue := fmt.Sprintf("results_%s", reg.ID)

		// Utilisation d'un canal temporaire pour la d√©claration (DevOps Best Practice)
		tmpCh, err := o.Conn.Channel()
		if err != nil {
			log.Printf("[RMQ] Erreur canal temporaire: %v", err)
			continue
		}

		_, err = tmpCh.QueueDeclare(
			resQueue,
			true, // Durable obligatoire pour les Streams
			false,
			false,
			false,
			amqp.Table{"x-queue-type": "stream", "x-expires": 60000}, // 7 jours en ms
		)
		tmpCh.Close()

		if err == nil {
			o.safeGo(func() {
				o.ConsumeWorkerResults(resQueue)
			})
			o.BroadcastToUI("WORKER_JOIN", reg) //
			log.Printf("[SYNC] Worker %s (%s) pr√™t.", reg.ID, reg.Language)
		}
	}
}

```

## orchestrator/internal/orchestrator/types.go

```go
package orchestrator

type WorkerRegistration struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	// Utilise "-" pour que ce ne soit pas envoy√© √† l'IHM, mais reste accessible en Go
	LastSeen int64 `json:"_"`
}

type AdminTask struct {
	TaskID  string                 `json:"task_id"`
	Handler string                 `json:"handler"`
	StartAt int64                  `json:"start_at"`
	Params  map[string]interface{} `json:"params"`
}

type WorkerResult struct {
	WorkerID  string      `json:"worker_id"`
	TaskID    string      `json:"task_id"`
	Handler   string      `json:"handler"`
	Index     int         `json:"index"`
	Metadata  interface{} `json:"metadata"`
	Timestamp int64       `json:"timestamp"`
}

```

## orchestrator/tmp/build-errors.log

```text
exit status 1exit status 1exit status 1exit status 1exit status 1exit status 1exit status 1exit status 1exit status 1
```

## orchestrator/tmp/main

> Fichier binaire non inclus (29445875 octets)

## project_export.log

```text
[2025-12-25 19:28:03] Source  : .
[2025-12-25 19:28:03] Sortie  : project_export.md
[2025-12-25 19:28:03] Fichiers trouv√©s (avant filtre): 181
[2025-12-25 19:28:03] Fichiers √† concat√©ner (apr√®s filtre): 180 (exclus auto:1 dir:0 file:0)
[2025-12-25 19:28:03] Concat√®ne [1] .env (size=627)
[2025-12-25 19:28:03] Concat√®ne [2] .git/COMMIT_EDITMSG (size=104)
[2025-12-25 19:28:03] Concat√®ne [3] .git/FETCH_HEAD (size=87)
[2025-12-25 19:28:03] Concat√®ne [4] .git/HEAD (size=23)
[2025-12-25 19:28:03] Concat√®ne [5] .git/ORIG_HEAD (size=41)
[2025-12-25 19:28:03] Concat√®ne [6] .git/config (size=348)
[2025-12-25 19:28:03] Concat√®ne [7] .git/description (size=73)
[2025-12-25 19:28:03] Concat√®ne [8] .git/hooks/applypatch-msg.sample (size=478)
[2025-12-25 19:28:03] Concat√®ne [9] .git/hooks/commit-msg.sample (size=896)
[2025-12-25 19:28:03] Concat√®ne [10] .git/hooks/fsmonitor-watchman.sample (size=4726)
[2025-12-25 19:28:03] Concat√®ne [11] .git/hooks/post-update.sample (size=189)
[2025-12-25 19:28:03] Concat√®ne [12] .git/hooks/pre-applypatch.sample (size=424)
[2025-12-25 19:28:03] Concat√®ne [13] .git/hooks/pre-commit.sample (size=1643)
[2025-12-25 19:28:03] Concat√®ne [14] .git/hooks/pre-merge-commit.sample (size=416)
[2025-12-25 19:28:03] Concat√®ne [15] .git/hooks/pre-push.sample (size=1374)
[2025-12-25 19:28:03] Concat√®ne [16] .git/hooks/pre-rebase.sample (size=4898)
[2025-12-25 19:28:03] Concat√®ne [17] .git/hooks/pre-receive.sample (size=544)
[2025-12-25 19:28:03] Concat√®ne [18] .git/hooks/prepare-commit-msg.sample (size=1492)
[2025-12-25 19:28:03] Concat√®ne [19] .git/hooks/push-to-checkout.sample (size=2783)
[2025-12-25 19:28:03] Concat√®ne [20] .git/hooks/sendemail-validate.sample (size=2308)
[2025-12-25 19:28:03] Concat√®ne [21] .git/hooks/update.sample (size=3650)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/index ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [23] .git/info/exclude (size=240)
[2025-12-25 19:28:03] Concat√®ne [24] .git/logs/HEAD (size=1271)
[2025-12-25 19:28:03] Concat√®ne [25] .git/logs/refs/heads/master (size=1271)
[2025-12-25 19:28:03] Concat√®ne [26] .git/logs/refs/remotes/origin/master (size=576)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/07/ae975db35e66033dec0b124a5beb703025248f ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/0c/4fc5ef3fa24178bd7db7c90efb69907387eee0 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/0c/636d86bffac4f802278a4a7178d9a27193f0ec ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/0e/ee988462b970f9b69b1449ac6c54ca6a9a5754 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/11/c0a428ae63867970c057e4da81ea673a9b10a4 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/12/a122ea7678e77dcb0ffc86aeca9c7d69bd816d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/12/d88b2843b5c4a93c94343f51967ec3b4538f6d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/16/f77612eccd94213863c5ad7ebd2e9fe8dd009b ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/1a/fb529027bfd3b244eb797399e6c0966eecbeb5 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/1d/17b6c4d2312673094dc22ccfeac6262f182d81 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/1e/9e3512be43ceac03f749b3e1bd81a4e0899e7f ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [38] .git/objects/20/708f6a2ef877be2ccc2101160dbce8791b4604 (size=273)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/22/1181d2632a36a783df61e3c0345b3ede421068 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [40] .git/objects/22/f7a387d14bb1ff975b8f8a909c2760552ab6b2 (size=577)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/24/3e0c81c7ad72a6d3e8acb2dec9be096a1193da ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/26/1363925bedc249329ef083cf913018d63c62f5 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/33/c0d87c6d62f2f9d9a4e40f8f25ed90e4056ba3 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/34/2ec37e7c6d81a03252d768ba5b12b3f726e491 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/37/ae47ae280e513c8aa2129bd8132aefba977851 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/3d/2251a550284deb468cc7e650b8b3aee7127908 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/3e/2093992b9bec6d3f7c01739e68ac39c300d3d1 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [48] .git/objects/3e/20ed345bffc4e3750b64a9529fa7d180030b91 (size=275)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/3f/a8c8c5920b5fd88bbfd47984588cab3f7465c5 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/42/4bb4fbbc66cf6af15cf54a2ff469815c3faddf ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/47/1c30f56909f444ef7f38a7049ab5233f1de6fc ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/49/6c619e5bce5b8e4b24d430a73bbd67e2f4f31e ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/49/a6f4e3b2dce8aee8a067bf9e18b853b66fd19d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [54] .git/objects/4b/37f9a11f378cc7e66ce3371d0508022174dd30 (size=237)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/4c/aeed4572829ae912ba1a0d88c87af1bc82e713 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/51/892cf7e7cd672d3b5d6d471bbb562d04a8e164 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/5a/50d4c26668c2b5b465756529b5868e7f109e11 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [58] .git/objects/5a/8af5557d17aa3d37ee1e7841508587a0d49f19 (size=228)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/5a/e095b7d3a678a69e6808cb3856ccd3ded0bdf5 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/5d/175e2933d1e59016f42d4e2f2441db2fef311d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [61] .git/objects/5d/ab255ae1904e7a27735e3c9641a108451695d3 (size=206)
[2025-12-25 19:28:03] Concat√®ne [62] .git/objects/5e/1af68ac259b3f3f34291e2803d602d13ce9fd5 (size=206)
[2025-12-25 19:28:03] Concat√®ne [63] .git/objects/5f/7131495b2d34d83cb08d08a31f739f2fa24949 (size=475)
[2025-12-25 19:28:03] Concat√®ne [64] .git/objects/5f/b088a58ba421fe022337f7079466ed09d22bb1 (size=201)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/61/9c750eb7e119698413b3f6346eb8c14f603684 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/64/72bab105a7266bf1324976307cfdc40c9fba6d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/67/057bc258203607838b49236920023738fec519 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [68] .git/objects/67/733188b7153a7ce368444a9509968142172eea (size=406)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/69/debbb92bbcb24a4513e4cf15034b3f40c2f338 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/6c/cc48a05421e0f3e40a0cfd4bc069663295c1c7 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/6d/6408a5bb8411448f2ee0f3e05078f8e13ff684 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/6f/9367e4ddeff2ceee6e1b6aa24bbf00f029cb75 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/71/d1529cacb1c8b66dff4899a11595712c0ff2ca ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/73/a6ec090865db226e94fa131701d00746aaa7a9 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/75/cdea54465f53768e9edc7c31fa5a226602af79 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/75/cfbdfbadfe07187291e5446058a8c6349db908 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [77] .git/objects/78/75ca23928e02658b93e6524a8c99458e7c39e6 (size=235)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/79/86982687d75da66eb87f7a2588dd8d33f06898 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7b/6628ade4bd1e5bf8024768ac26ff688ef83082 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7c/d49f16753e2d5e87e82317cc62ba5ad25fcd1a ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7d/14ca0ceb65a100b65da50c13c4ed00423b7abd ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7d/d15a1c616cdd4892d04ffee27e0b1f8e0240bb ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7e/8da60904b656207a5010b7027b47e651c85c0c ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/7f/caa3a05cd27baa703ea6e7989472bdc7f6e469 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/81/692525cea495f555182927f8f55584d98020f9 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/83/fc2d1afb8f6583a9ac47503781a5f53b7e41c8 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/84/058f295333b603c0d827b8f70f6f199e64a3df ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/84/6efe57e7f1aba3621ccda2c424a11089a02b9b ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [89] .git/objects/84/dc0203f93f52853e23371eb3c46d49449ccf3e (size=214)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/87/5968d8f81678a6e9a858eebe90a2b957759531 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/87/63c1d9aa8a5da1da10bdcac57f1fc5ebb09491 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [92] .git/objects/89/13571fad3aed041bd6f5a6fbb57446cf28ec84 (size=219)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/8a/f4c7ac97d6f232d2b6a4a1e6041d7ccb6ac8bb ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/8e/cc24ea2efe2596320b847a326edc2296cd2c7f ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/91/f5acd81ec45ecd315e27d23991247bd92e1979 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/93/15641dac3ec2c85995ad959b95037cc83827c5 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/96/47bb3ee733f230e94decddebc7276f600e2b32 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/9d/c71497ae68a4e83a6d720c3ad2fe6abd38d86f ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/9d/cf46457a2239e54fbe4a9f00f2826996522e04 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/a2/91f3e344cf5c4aef3053ba2cb40e91b5fc39b3 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/a5/b9c675de4175687d7fff51d0dab61507fb118d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [102] .git/objects/ab/ece7364f2ed1cdc9b439c4c5374ab3d9dfd837 (size=286)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/ac/d9d78e43721d13834d5d0f2c10a60001a496b6 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/ad/7005ae32ba2e370bdae2663559142ddea29a5b ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/af/3573daae845bf3f07ef5299640f9c6f5af2201 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/b0/aaf9bab41421336a5b566acb37ad38096d48d3 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/b1/9ef482e0ef1f6fb118834ce1c56aec63821e1d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/b5/3633d98bf09f65d7ffad41851e7dc00426fc7a ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/b7/6e4ffecf8364809c61229d2644b578f6a94157 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [110] .git/objects/b7/9c49648d6cba60ded7e23b9eccc64959b1031c (size=457)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/b7/e1fe52398f8f442711a0d0470a6318c0f10df3 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/be/0c220c16ac5a2391490b3f9618e5f9454345aa ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/bf/2fb9f446bf88683793ffeb730ffe33654b8be1 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [114] .git/objects/c2/7fb11e9c9c3fe264ad1fdf894bbb8214c1d6a7 (size=216)
[2025-12-25 19:28:03] Concat√®ne [115] .git/objects/c4/e58bd3837ab011a7decff3a984391061896366 (size=294)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/c4/f763aaa6fea4a56f4c36a8a0bd290b983f720b ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/c7/9719e52ab405c026c4d3c48899017ff6545063 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/c8/7dd4961de5d2b0546460efbf71169d6fbde2c0 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/c9/a83108423860e7c4e872bdb748dfb3d2f60be9 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/ca/702652572c3cdc89b07f30a4f264c3cb78c40f ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [121] .git/objects/cb/01ac9c40ec00fd41ff670bbc86735946ff34a9 (size=205)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/d4/55b2fc342a290696ff081184fe4d5f18c87b7c ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/d4/ad699edca1759fdd98c08d197aea6fe8e7e202 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/d7/0c9ad928febd9287bcb9079f5857db1af09c06 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/d7/b6bdacc7db9681e4d5c83bbce813d0dae49244 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/dc/dd5be855c5c415d5179b9375b819fe9c3daf06 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/df/c5b6c588b2753fb511e78f198853374d94ddbc ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/e1/13dd3b75bdcf3d891223e4ee5a8e376a739d06 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [129] .git/objects/e2/0e09349e91e8cafe04f0cd25d2b449b3459170 (size=153)
[2025-12-25 19:28:03] Concat√®ne [130] .git/objects/e2/b928d9112e1dc2bf98d9006e9a0bf2fc8cbf93 (size=470)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/e3/8fb63d9a6cac854ef98242004daeaa4ab8d694 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/e4/0d4a3fbc0649821696f81aeb021bb21fbda101 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/eb/3dda2764d74686f269b5551c818e3a4028218d ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/ec/d75440765f8c854fc0dd750e7f3249b9194b94 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/ee/d6e7bc041cbff29407b9c98dd75b168a5869fa ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [136] .git/objects/ef/327c55831b06abfc2b987608b12db85728ffbf (size=365)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/f1/2d44b84ecbccc4b653c7fc833ad2dd7373db46 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/f5/c1c09205f17763c8289c63c6d35fd4df64c6a8 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/f6/64bd3652c06f1111cd71c7133c9409dd882479 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [140] .git/objects/f8/786ea322dc83337b57a86a0e86612465bbd5b5 (size=363)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : .git/objects/fa/10cc7e0d046698a0c2b77544ed2a7167a49d25 ‚Äî r√©f√©renc√© mais non inclus
[2025-12-25 19:28:03] Concat√®ne [142] .git/refs/heads/master (size=41)
[2025-12-25 19:28:03] Concat√®ne [143] .git/refs/remotes/origin/master (size=41)
[2025-12-25 19:28:03] Concat√®ne [144] .github/copilot-instructions.md (size=5199)
[2025-12-25 19:28:03] Concat√®ne [145] .gitignore (size=123)
[2025-12-25 19:28:03] Concat√®ne [146] compose.yml (size=2379)
[2025-12-25 19:28:03] Concat√®ne [147] orchestrator/.air.toml (size=212)
[2025-12-25 19:28:03] Concat√®ne [148] orchestrator/Dockerfile (size=697)
[2025-12-25 19:28:03] Concat√®ne [149] orchestrator/cmd/server/main.go (size=1974)
[2025-12-25 19:28:03] Concat√®ne [150] orchestrator/go.mod (size=1666)
[2025-12-25 19:28:03] Concat√®ne [151] orchestrator/go.sum (size=8015)
[2025-12-25 19:28:03] Concat√®ne [152] orchestrator/internal/api/hub.go (size=904)
[2025-12-25 19:28:03] Concat√®ne [153] orchestrator/internal/api/routes.go (size=1566)
[2025-12-25 19:28:03] Concat√®ne [154] orchestrator/internal/api/web/index.html (size=2206)
[2025-12-25 19:28:03] Concat√®ne [155] orchestrator/internal/api/web/js/app.js (size=1070)
[2025-12-25 19:28:03] Concat√®ne [156] orchestrator/internal/api/web/js/charts.js (size=1226)
[2025-12-25 19:28:03] Concat√®ne [157] orchestrator/internal/api/web/js/ui.js (size=4425)
[2025-12-25 19:28:03] Concat√®ne [158] orchestrator/internal/orchestrator/engine.go (size=3532)
[2025-12-25 19:28:03] Concat√®ne [159] orchestrator/internal/orchestrator/rabbitmq.go (size=4823)
[2025-12-25 19:28:03] Concat√®ne [160] orchestrator/internal/orchestrator/types.go (size=730)
[2025-12-25 19:28:03] Concat√®ne [161] orchestrator/tmp/build-errors.log (size=117)
[2025-12-25 19:28:03] ‚ÑπÔ∏è  Binaire : orchestrator/tmp/main ‚Äî r√©f√©renc√© mais non inclus

```

## worker-go/Dockerfile

```text
# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache protoc protobuf-dev git ca-certificates

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# 2. G√©n√©ration du code gRPC
COPY proto/ ./proto/
# On cr√©e le dossier de destination
RUN mkdir -p internal/pb 

# Commande protoc align√©e sur le module "fibo-worker" 
RUN protoc --proto_path=./proto \
    --go_out=. --go_opt=module=fibo-worker \
    --go-grpc_out=. --go-grpc_opt=module=fibo-worker \
    ./proto/sync.proto

# 3. Copie du code source et Build
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/worker ./cmd/worker/main.go

# --- Stage 2: Image Finale (Scratch pour ton cluster k3s) ---
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
COPY --from=builder /app/worker .

EXPOSE 8081 50051
ENTRYPOINT ["./worker"]
```

## worker-go/cmd/worker/main.go

```go
package main

import (
	"fibo-worker/internal/api"
	"fibo-worker/internal/worker"
	"os"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://bench_go:9q8s7d9qs87dqs654dq6s54d6qs54dqs321dqs2d1qs98d7qs9d8q7@192.168.1.12:5672/benchmarks"
	}

	// 1. Lancement du moteur Worker (Async)
	engine := worker.NewEngine(amqpURL)
	go engine.Start()

	// 2. Lancement de l'API de sant√© et IO (Sync)
	r := api.SetupRouter()
	r.Run(":8081")
}

```

## worker-go/go.mod

```text
module fibo-worker

go 1.25.1

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/streadway/amqp v1.1.0
)

require (
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

```

## worker-go/go.sum

```text
github.com/bytedance/sonic v1.14.0 h1:/OfKt8HFw0kh2rj8N0F6C/qPGRESq0BbaNZgcNXXzQQ=
github.com/bytedance/sonic v1.14.0/go.mod h1:WoEbx8WTcFJfzCe0hbmyTGrfjt8PzNEBdxlNUO24NhA=
github.com/bytedance/sonic/loader v0.3.0 h1:dskwH8edlzNMctoruo8FPTJDF3vLtDT0sXZwvZJyqeA=
github.com/bytedance/sonic/loader v0.3.0/go.mod h1:N8A3vUdtUebEY2/VQC0MyhYeKUFosQU6FxH2JmUe6VI=
github.com/cloudwego/base64x v0.1.6 h1:t11wG9AECkCDk5fMSoxmufanudBtJ+/HemLstXDLI2M=
github.com/cloudwego/base64x v0.1.6/go.mod h1:OFcloc187FXDaYHvrNIjxSe8ncn0OOM8gEHfghB2IPU=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/gabriel-vasile/mimetype v1.4.8 h1:FfZ3gj38NjllZIeJAmMhr+qKL8Wu+nOoI3GqacKw1NM=
github.com/gabriel-vasile/mimetype v1.4.8/go.mod h1:ByKUIKGjh1ODkGM1asKUbQZOLGrPjydw3hYPU2YU9t8=
github.com/gin-contrib/sse v1.1.0 h1:n0w2GMuUpWDVp7qSpvze6fAu9iRxJY4Hmj6AmBOU05w=
github.com/gin-contrib/sse v1.1.0/go.mod h1:hxRZ5gVpWMT7Z0B0gSNYqqsSCNIJMjzvm6fqCz9vjwM=
github.com/gin-gonic/gin v1.11.0 h1:OW/6PLjyusp2PPXtyxKHU0RbX6I/l28FTdDlae5ueWk=
github.com/gin-gonic/gin v1.11.0/go.mod h1:+iq/FyxlGzII0KHiBGjuNn4UNENUlKbGlNmc+W50Dls=
github.com/go-playground/assert/v2 v2.2.0 h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=
github.com/go-playground/assert/v2 v2.2.0/go.mod h1:VDjEfimB/XKnb+ZQfWdccd7VUvScMdVu0Titje2rxJ4=
github.com/go-playground/locales v0.14.1 h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=
github.com/go-playground/locales v0.14.1/go.mod h1:hxrqLVvrK65+Rwrd5Fc6F2O76J/NuW9t0sjnWqG1slY=
github.com/go-playground/universal-translator v0.18.1 h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=
github.com/go-playground/universal-translator v0.18.1/go.mod h1:xekY+UJKNuX9WP91TpwSH2VMlDf28Uj24BCp08ZFTUY=
github.com/go-playground/validator/v10 v10.27.0 h1:w8+XrWVMhGkxOaaowyKH35gFydVHOvC0/uWoy2Fzwn4=
github.com/go-playground/validator/v10 v10.27.0/go.mod h1:I5QpIEbmr8On7W0TktmJAumgzX4CA1XNl4ZmDuVHKKo=
github.com/goccy/go-json v0.10.2 h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=
github.com/goccy/go-json v0.10.2/go.mod h1:6MelG93GURQebXPDq3khkgXZkazVtN9CRI+MGFi0w8I=
github.com/goccy/go-yaml v1.18.0 h1:8W7wMFS12Pcas7KU+VVkaiCng+kG8QiFeFwzFb+rwuw=
github.com/goccy/go-yaml v1.18.0/go.mod h1:XBurs7gK8ATbW4ZPGKgcbrY1Br56PdM69F7LkFRi1kA=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/json-iterator/go v1.1.12 h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=
github.com/json-iterator/go v1.1.12/go.mod h1:e30LSqwooZae/UwlEbR2852Gd8hjQvJoHmT4TnhNGBo=
github.com/klauspost/cpuid/v2 v2.3.0 h1:S4CRMLnYUhGeDFDqkGriYKdfoFlDnMtqTiI/sFzhA9Y=
github.com/klauspost/cpuid/v2 v2.3.0/go.mod h1:hqwkgyIinND0mEev00jJYCxPNVRVXFQeu1XKlok6oO0=
github.com/leodido/go-urn v1.4.0 h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=
github.com/leodido/go-urn v1.4.0/go.mod h1:bvxc+MVxLKB4z00jd1z+Dvzr47oO32F/QSNjSBOlFxI=
github.com/mattn/go-isatty v0.0.20 h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=
github.com/mattn/go-isatty v0.0.20/go.mod h1:W+V8PltTTMOvKvAeJH7IuucS94S2C6jfK/D7dTCTo3Y=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 h1:ZqeYNhU3OHLH3mGKHDcjJRFFRrJa6eAM5H+CtDdOsPc=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v1.0.2 h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=
github.com/modern-go/reflect2 v1.0.2/go.mod h1:yWuevngMOJpCy52FWWMvUC8ws7m/LJsjYzDa0/r8luk=
github.com/pelletier/go-toml/v2 v2.2.4 h1:mye9XuhQ6gvn5h28+VilKrrPoQVanw5PMw/TB0t5Ec4=
github.com/pelletier/go-toml/v2 v2.2.4/go.mod h1:2gIqNv+qfxSVS7cM2xJQKtLSTLUE9V8t9Stt+h56mCY=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/quic-go/qpack v0.5.1 h1:giqksBPnT/HDtZ6VhtFKgoLOWmlyo9Ei6u9PqzIMbhI=
github.com/quic-go/qpack v0.5.1/go.mod h1:+PC4XFrEskIVkcLzpEkbLqq1uCoxPhQuvK5rH1ZgaEg=
github.com/quic-go/quic-go v0.54.0 h1:6s1YB9QotYI6Ospeiguknbp2Znb/jZYjZLRXn9kMQBg=
github.com/quic-go/quic-go v0.54.0/go.mod h1:e68ZEaCdyviluZmy44P6Iey98v/Wfz6HCjQEm+l8zTY=
github.com/streadway/amqp v1.1.0 h1:py12iX8XSyI7aN/3dUT8DFIDJazNJsVJdxNVEpnQTZM=
github.com/streadway/amqp v1.1.0/go.mod h1:WYSrTEYHOXHd0nwFeUXAe2G2hRnQT+deZJJf88uS9Bg=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.4.0/go.mod h1:YvHI0jy2hoMjB+UWwv71VJQ9isScKT/TqJzVSSt89Yw=
github.com/stretchr/objx v0.5.0/go.mod h1:Yh+to48EsGEfYuaHDzXPcE3xhTkx73EhmCGUpEOglKo=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.7.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.8.0/go.mod h1:yNjHg4UonilssWZ8iaSj1OCr/vHnekPRkoO+kdMU+MU=
github.com/stretchr/testify v1.8.1/go.mod h1:w2LPCIKwWwSfY2zedu0+kehJoqGctiVI29o6fzry7u4=
github.com/stretchr/testify v1.11.1 h1:7s2iGBzp5EwR7/aIZr8ao5+dra3wiQyKjjFuvgVKu7U=
github.com/stretchr/testify v1.11.1/go.mod h1:wZwfW3scLgRK+23gO65QZefKpKQRnfz6sD981Nm4B6U=
github.com/twitchyliquid64/golang-asm v0.15.1 h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=
github.com/twitchyliquid64/golang-asm v0.15.1/go.mod h1:a1lVb/DtPvCB8fslRZhAngC2+aY1QWCk3Cedj/Gdt08=
github.com/ugorji/go/codec v1.3.0 h1:Qd2W2sQawAfG8XSvzwhBeoGq71zXOC/Q1E9y/wUcsUA=
github.com/ugorji/go/codec v1.3.0/go.mod h1:pRBVtBSKl77K30Bv8R2P+cLSGaTtex6fsA2Wjqmfxj4=
go.uber.org/mock v0.5.0 h1:KAMbZvZPyBPWgD14IrIQ38QCyjwpvVVV6K/bHl1IwQU=
go.uber.org/mock v0.5.0/go.mod h1:ge71pBPLYDk7QIi1LupWxdAykm7KIEFchiOqd6z7qMM=
golang.org/x/arch v0.20.0 h1:dx1zTU0MAE98U+TQ8BLl7XsJbgze2WnNKF/8tGp/Q6c=
golang.org/x/arch v0.20.0/go.mod h1:bdwinDaKcfZUGpH09BB7ZmOfhalA8lQdzl62l8gGWsk=
golang.org/x/crypto v0.40.0 h1:r4x+VvoG5Fm+eJcxMaY8CQM7Lb0l1lsmjGBQ6s8BfKM=
golang.org/x/crypto v0.40.0/go.mod h1:Qr1vMER5WyS2dfPHAlsOj01wgLbsyWtFn/aY+5+ZdxY=
golang.org/x/mod v0.25.0 h1:n7a+ZbQKQA/Ysbyb0/6IbB1H/X41mKgbhfv7AfG/44w=
golang.org/x/mod v0.25.0/go.mod h1:IXM97Txy2VM4PJ3gI61r1YEk/gAj6zAHN3AdZt6S9Ww=
golang.org/x/net v0.42.0 h1:jzkYrhi3YQWD6MLBJcsklgQsoAcw89EcZbJw8Z614hs=
golang.org/x/net v0.42.0/go.mod h1:FF1RA5d3u7nAYA4z2TkclSCKh68eSXtiFwcWQpPXdt8=
golang.org/x/sync v0.16.0 h1:ycBJEhp9p4vXvUZNszeOq0kGTPghopOL8q0fq3vstxw=
golang.org/x/sync v0.16.0/go.mod h1:1dzgHSNfp02xaA81J2MS99Qcpr2w7fw1gpm99rleRqA=
golang.org/x/sys v0.6.0/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
golang.org/x/sys v0.35.0 h1:vz1N37gP5bs89s7He8XuIYXpyY0+QlsKmzipCbUtyxI=
golang.org/x/sys v0.35.0/go.mod h1:BJP2sWEmIv4KK5OTEluFJCKSidICx8ciO85XgH3Ak8k=
golang.org/x/text v0.27.0 h1:4fGWRpyh641NLlecmyl4LOe6yDdfaYNrGb2zdfo4JV4=
golang.org/x/text v0.27.0/go.mod h1:1D28KMCvyooCX9hBiosv5Tz/+YLxj0j7XhWjpSUF7CU=
golang.org/x/tools v0.34.0 h1:qIpSLOxeCYGg9TrcJokLBG4KFA6d795g0xkBkiESGlo=
golang.org/x/tools v0.34.0/go.mod h1:pAP9OwEaY1CAW3HOmg3hLZC5Z0CCmzjAF2UQMSqNARg=
google.golang.org/protobuf v1.36.9 h1:w2gp2mA27hUeUzj9Ex9FBjsBm40zfaDtEWow293U7Iw=
google.golang.org/protobuf v1.36.9/go.mod h1:fuxRtAxBytpl4zzqUh6/eyUujkJdNiuEkXntxiD/uRU=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=

```

## worker-go/internal/api/router.go

```go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Futur handler pour les tests de transfert de fichiers
	r.POST("/io-test", func(c *gin.Context) {
		// Logique de r√©ception de fichier pour benchmark
		c.Status(http.StatusAccepted)
	})

	return r
}

```

## worker-go/internal/worker/engine.go

```go
package worker

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Engine struct {
	ID      string
	Conn    *amqp.Connection
	Channel *amqp.Channel
	AMQPURL string
	Mu      sync.Mutex
}

func NewEngine(url string) *Engine {
	return &Engine{
		ID:      GenerateID(),
		AMQPURL: url,
	}
}

func (e *Engine) Start() {
	for {
		log.Printf("[WORKER] Tentative de connexion RMQ sur %s", e.AMQPURL)
		conn, err := amqp.Dial(e.AMQPURL)
		if err != nil {
			log.Printf("[ERROR] √âchec connexion: %v. Re-tentative dans 15s...", err)
			time.Sleep(15 * time.Second)
			continue
		}

		e.Conn = conn
		ch, err := conn.Channel()
		if err != nil {
			log.Printf("[ERROR] Impossible d'ouvrir un canal: %v", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		e.Mu.Lock()
		e.Channel = ch
		e.Mu.Unlock()

		// Configuration de l'infrastructure minimale pour le worker
		if err := e.setupInfra(); err != nil {
			log.Printf("[ERROR] Erreur setup infra: %v", err)
			conn.Close()
			time.Sleep(15 * time.Second)
			continue
		}

		// Phase 1 : Signalement de pr√©sence
		e.register()
		// D√©marrage du heartbeat
		go e.startHeartbeat()

		// Phase 2 : √âcoute des ordres (Exchange Fanout)
		go e.listenTasks()

		// Surveillance de la sant√© de la connexion
		closeChan := make(chan *amqp.Error)
		e.Conn.NotifyClose(closeChan)

		log.Println("[WORKER] Connect√© et pr√™t √† recevoir des t√¢ches.")

		err = <-closeChan
		log.Printf("[WARN] Connexion perdue: %v. Relance de la boucle...", err)
	}
}

func (e *Engine) setupInfra() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	// D√©claration de la queue de pr√©sence (doit √™tre identique √† l'orchestrateur)
	_, err := e.Channel.QueueDeclare(
		"isReady",
		true, // Durable: true (match orchestrateur)
		false,
		false,
		false,
		nil,
	)
	return err
}

func (e *Engine) register() {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	reg, _ := json.Marshal(WorkerRegistration{ID: e.ID, Language: "go"})
	err := e.Channel.Publish(
		"",
		"isReady",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        reg,
		},
	)
	if err != nil {
		log.Printf("[ERROR] √âchec de l'enregistrement: %v", err)
	} else {
		log.Printf("[PHASE 1] Worker enregistr√© avec l'ID: %s", e.ID)
	}
}

func (e *Engine) startHeartbeat() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		health := map[string]interface{}{
			"worker_id": e.ID,
			"status":    "active",
			"ram":       m.Alloc / 1024 / 1024,
			"cpu":       runtime.NumGoroutine(),
			"timestamp": time.Now().Unix(),
			// Simuler I/O pour l'instant (√† lier √† /proc/net/dev plus tard)
			"net_io":  "2.4MB/s",
			"disk_io": "150KB/s",
		}

		body, _ := json.Marshal(health)
		e.Mu.Lock()
		if e.Channel != nil {
			e.Channel.Publish("", "worker_health", false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		}
		e.Mu.Unlock()
	}
}

```

## worker-go/internal/worker/handlers.go

```go
package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"runtime"
	"time"

	"github.com/streadway/amqp"
)

func (e *Engine) listenTasks() {
	e.Mu.Lock()
	ch := e.Channel
	e.Mu.Unlock()

	// D√©claration de l'exchange de diffusion
	ch.ExchangeDeclare("fibo_admin_exchange", "fanout", true, false, false, false, nil)

	// Cr√©ation d'une queue temporaire pour ce worker sp√©cifique
	q, _ := ch.QueueDeclare("", false, true, true, false, nil)
	ch.QueueBind(q.Name, "", "fibo_admin_exchange", false, nil)

	msgs, _ := ch.Consume(q.Name, "", true, false, false, false, nil)

	for d := range msgs {
		var task AdminTask
		if err := json.Unmarshal(d.Body, &task); err != nil {
			continue
		}

		// Synchronisation temporelle (Phase 3)
		wait := time.Until(time.Unix(task.StartAt, 0))
		if wait > 0 {
			log.Printf("[SYNC] Attente de %v pour d√©marrage synchrone...", wait)
			time.Sleep(wait)
		}

		log.Printf("[TASK] D√©marrage du calcul: %s (ID: %s)", task.Handler, task.TaskID)
		e.runHandler(task)

		// Nettoyage m√©moire apr√®s calcul intensif pour ton cluster IA
		runtime.GC()
		log.Println("[CLEAN] Ressources lib√©r√©es apr√®s t√¢che.")
	}
}

func (e *Engine) runHandler(task AdminTask) {
	switch task.Handler {
	case "fibonacci":
		e.handleFibo(task)
	default:
		log.Printf("[WARN] Handler inconnu: %s", task.Handler)
	}
}

func (e *Engine) handleFibo(task AdminTask) {
	limit := 400000
	if val, ok := task.Params["limit"].(float64); ok {
		limit = int(val)
	}

	a, b := big.NewInt(0), big.NewInt(1)
	resQueue := "results_" + e.ID

	// On pr√©pare les stats syst√®me une seule fois ou p√©riodiquement
	var m runtime.MemStats

	for i := 0; i <= limit; i++ {
		a.Add(a, b)
		a, b = b, a

		// OPTIMISATION : On ne lit la RAM que toutes les 1000 it√©rations
		// pour ne pas ralentir le calcul pur, mais on ENVOIE chaque message.
		if i%1000 == 0 {
			runtime.ReadMemStats(&m)
		}

		res := WorkerResult{
			WorkerID:  e.ID,
			TaskID:    task.TaskID,
			Handler:   "fibonacci",
			Index:     i,
			Timestamp: time.Now().UnixMilli(),
			Metadata: map[string]interface{}{
				"cpu":     runtime.NumGoroutine(),
				"ram":     m.Alloc / 1024 / 1024,
				"net_io":  fmt.Sprintf("%.2f KB", float64(len(a.Bits())*8)/1024.0),
				"disk_io": "0.1 MB/s",
			},
		}

		body, _ := json.Marshal(res)

		// PUBLICATION √Ä CHAQUE IT√âRATION
		e.Channel.Publish(
			"",
			resQueue,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
	}
	log.Printf("[FINISH] %d messages envoy√©s.", limit)
}

```

## worker-go/internal/worker/types.go

```go
package worker

import (
	"crypto/sha256"
	"fmt"
	"os"
)

type WorkerRegistration struct {
	ID       string `json:"id"`
	Language string `json:"language"`
}

type AdminTask struct {
	TaskID  string                 `json:"task_id"`
	Handler string                 `json:"handler"`
	StartAt int64                  `json:"start_at"`
	Params  map[string]interface{} `json:"params"`
}

type WorkerResult struct {
	WorkerID  string      `json:"worker_id"`
	TaskID    string      `json:"task_id"`
	Handler   string      `json:"handler"`
	Index     int         `json:"index"`
	Metadata  interface{} `json:"metadata"`
	Timestamp int64       `json:"timestamp"`
}

// GenerateID cr√©e l'identifiant unique SHA256 du worker au d√©marrage
func GenerateID() string {
	hostname, _ := os.Hostname()
	data := fmt.Sprintf("%s-%d", hostname, os.Getpid())
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

```

## worker-go/proto/sync.proto

```text
syntax = "proto3";

package sync;

// CORRECTION : Le chemin DOIT commencer par le nom du module d√©fini dans go.mod
option go_package = "fibo-worker/internal/pb"; 

service Barrier {
    rpc WaitToStart (Empty) returns (StartSignal);
}

message Empty {}
message StartSignal {
    string message = 1;
}
```

## worker-node/Dockerfile

```text
FROM node:20-alpine

WORKDIR /app

RUN npm install @grpc/grpc-js @grpc/proto-loader

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# --- CORRECTION ICI ---
# On sp√©cifie le dossier source 'node/'
COPY node/index.js .

CMD ["node", "index.js"]
```

## worker-node/index.js

```javascript
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const packageDefinition = protoLoader.loadSync('sync.proto');
const syncProto = grpc.loadPackageDefinition(packageDefinition).sync;

function runFibo() {
    for (let run = 1; run <= 10; run++) {
        let a = 0n, b = 1n;
        const start = Date.now();
        for (let i = 0; i <= 400000; i++) {
            [a, b] = [b, a + b];
            if (i % 10000 === 0 && i > 0) {
                console.log(`[NODE] Run ${run} - ${i} iters - Temps: ${(Date.now() - start)/1000}s`);
            }
        }
    }
}

const client = new syncProto.Barrier('fibo-go:50051', grpc.credentials.createInsecure());
console.log("Node pr√™t, en attente du signal...");
client.waitToStart({}, (err) => {
    if (err) console.error(err);
    else runFibo();
});
```

## worker-python/Dockerfile

```text
FROM python:3.12-slim

WORKDIR /app

RUN pip install --no-cache-dir grpcio grpcio-tools

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# G√©n√©ration du code Python
RUN python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. sync.proto

# --- CORRECTION ICI ---
# On sp√©cifie le dossier source 'python/'
COPY python/main.py .

CMD ["python", "main.py"]
```

## worker-python/main.py

```python
import grpc
import time
from sync_pb2 import Empty
from sync_pb2_grpc import BarrierStub

def run_fibo():
    for run in range(1, 11):
        a, b = 0, 1
        start_time = time.time()
        for i in range(400001):
            a, b = b, a + b
            if i % 10000 == 0 and i > 0:
                print(f"[PYTHON] Run {run} - {i} iters - Temps: {time.time() - start_time:.4f}s")

if __name__ == "__main__":
    with grpc.insecure_channel('fibo-go:50051') as channel:
        stub = BarrierStub(channel)
        print("Python pr√™t, en attente du signal...")
        stub.WaitToStart(Empty())
        run_fibo()
```

## worker-rust/Cargo.toml

```toml
[package]
name = "fibo-rust"
version = "0.1.0"
edition = "2021"

[dependencies]
# RabbitMQ client
lapin = "2.3"
# Async runtime
tokio = { version = "1.0", features = ["full"] }
# Serialization
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
# Math & Utils
num-bigint = "0.4"
num-traits = "0.2"
futures-util = "0.3"
amqprs = "1.5" # Alternative possible, mais restons sur lapin pour la flexibilit√©

```

## worker-rust/Dockerfile

```text
FROM rust:1.83-alpine AS builder
RUN apk add --no-cache musl-dev protobuf-dev

WORKDIR /usr/src/app
COPY proto/ ../proto/
COPY rust/ .
RUN cargo build --release

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /usr/src/app/target/release/fibo-rust .
CMD ["./fibo-rust"]
```

## worker-rust/build.rs

```rust
fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::compile_protos("../proto/sync.proto")?;
    Ok(())
}
```

## worker-rust/src/main.rs

```rust
use futures_util::StreamExt;
use lapin::{options::*, types::FieldTable, BasicProperties, Connection, ConnectionProperties};
use num_bigint::BigInt;
use num_traits::{One, Zero};
use serde::{Deserialize, Serialize};
use std::time::Instant;

#[derive(Serialize, Deserialize, Debug)]
struct TaskMessage {
    language: String,
    serie: String,
    limit: u32,
}

#[derive(Serialize, Debug)]
struct ResultMessage {
    id: String,
    S√©rie: String,
    num: u32,
    value: String,
}

async fn run_fibo_and_publish(chan: &lapin::Channel, task: TaskMessage) {
    let mut a: BigInt = BigInt::zero();
    let mut b: BigInt = BigInt::one();
    let worker_id = "rust-worker-01".to_string();

    println!("[RUST] D√©marrage de la s√©rie : {}", task.serie);

    for i in 0..=task.limit {
        let temp = a.clone() + &b;
        a = b;
        b = temp;

        let res = ResultMessage {
            id: worker_id.clone(),
            S√©rie: task.serie.clone(),
            num: i,
            value: a.to_string(), // S√©rialisation BigInt en string
        };

        let payload = serde_json::to_vec(&res).unwrap();
        
        // Envoi imm√©diat √† RabbitMQ
        chan.basic_publish(
            "",
            "fibo_results",
            BasicPublishOptions::default(),
            &payload,
            BasicProperties::default(),
        )
        .await
        .expect("Erreur lors de la publication");

        if i % 10000 == 0 {
            println!("[RUST] {} it√©rations envoy√©es...", i);
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = std::env::var("AMQP_ADDR").unwrap_or_else(|_| "amqp://guest:guest@rabbitmq:5672/%2f".into());
    let conn = Connection::connect(&addr, ConnectionProperties::default()).await?;
    let channel = conn.create_channel().await?;

    // D√©claration des queues
    channel.queue_declare("fibo_tasks", QueueDeclareOptions::default(), FieldTable::default()).await?;
    channel.queue_declare("fibo_results", QueueDeclareOptions::default(), FieldTable::default()).await?;

    println!("Rust est connect√© √† RabbitMQ. En attente de t√¢ches...");

    let mut consumer = channel
        .basic_consume("fibo_tasks", "rust_consumer", BasicConsumeOptions::default(), FieldTable::default())
        .await?;

    while let Some(delivery) = consumer.next().await {
        let (_, delivery) = delivery.expect("error in consumer");
        let task: TaskMessage = serde_json::from_slice(&delivery.data)?;

        if task.language == "rust" || task.language == "all" {
            run_fibo_and_publish(&channel, task).await;
            channel.basic_ack(delivery.delivery_tag, BasicAckOptions::default()).await?;
        }
    }

    Ok(())
}
```

