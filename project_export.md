# Export de projet

_GÃ©nÃ©rÃ© le 2025-12-25T15:13:06+01:00_

## .git/COMMIT_EDITMSG

```text
feat: initial gRPC fibonacci benchmark with Go, Python and Node

```

## .git/HEAD

```text
ref: refs/heads/master

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

> Fichier binaire non inclus (1287 octets)

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

```

## .git/logs/refs/heads/master

```text
0000000000000000000000000000000000000000 9dcf46457a2239e54fbe4a9f00f2826996522e04 sicDANGBE <dansoug@gmail.com> 1766448047 +0100	commit (initial): feat: initial gRPC fibonacci benchmark with Go, Python and Node

```

## .git/logs/refs/remotes/origin/master

```text
0000000000000000000000000000000000000000 9dcf46457a2239e54fbe4a9f00f2826996522e04 sicDANGBE <dansoug@gmail.com> 1766448056 +0100	update by push

```

## .git/objects/12/a122ea7678e77dcb0ffc86aeca9c7d69bd816d

> Fichier binaire non inclus (848 octets)

## .git/objects/20/708f6a2ef877be2ccc2101160dbce8791b4604

```text
xUÁJÄ0†=ç)†zXÓ¢{POY•"»]ÂŠˆHÉ¦Ñ´™¤`©ÏÑ³m
â)“ùgæÿşcGXßŞ=ñb¶'4÷ëôú†úZ7„¼üe“sÈ„µ„ğ×XmAD]¥©ò¤h¥|;+5.ˆµ'äZ­ jÁ:MP?ARuIX±bæ;#Ó8—N»ÏCo†Ş‰ ÑL'$V
ö3f¤™K ÍìZÎq_ÍÓ‘0f*±Ówb,ÿõş,'CJ)°‚óGvÈ‹ä,ŸZ£Pğvè¥şÃÔcô^+['¬âÍlµ¤™¹²Fh“ÚÆ,l»$N%W,Ròù[‚€k
```

## .git/objects/37/ae47ae280e513c8aa2129bd8132aefba977851

> Fichier binaire non inclus (5707 octets)

## .git/objects/5a/50d4c26668c2b5b465756529b5868e7f109e11

> Fichier binaire non inclus (1105 octets)

## .git/objects/5d/ab255ae1904e7a27735e3c9641a108451695d3

```text
x}Án„0D{m¾"Çö€ÇèçQ$³VÛ¿o–=´e¥Ş,k<óÆıÌ½$Û¼|ò¸Ï^N±çä†!
X"˜Pˆì/{Ì^æ0{<»€sP!/ƒ¼"´-èİ›xı¥¹©ä7yÕP[ÀJƒÑ¦Ad¨n©«,èÊBO‘JÉ˜Æ’5l'õk½{P	úGµùÛEúYöÄîÓ’ycõ¨å–¸ª£N¹]=0ÁºL¶"‡m=}ØÎºîpv>lû}º†Jé¿Äïâq×r
```

## .git/objects/5e/1af68ac259b3f3f34291e2803d602d13ce9fd5

```text
xMO=KAµŞ_ñHsIu¨BH#XXG°”ıÊ1x·³Îî‰!ä¿;{pÁbŠ7ïkÆìğğütWÎ©Ú_°ÉÂ•7{c²õ_vˆPÎ+ì{¼¥@ßsD›1bœ,UR)"Ì˜8ÌÊv'rœ¬÷Ô!ØTš#p)]vá\‰ş\ƒ´şfê³kG”(?ä#^¬Hs^ ÙãC;ßùXµÛ×)×óë,Z´]ÖG’w{s5fŠ¥´W!.×ÛâŸr‰.U(XõÜ«ÿ—šc'
```

## .git/objects/5f/7131495b2d34d83cb08d08a31f739f2fa24949

```text
x”=nÜ0…SûSî±Lêoe)‚©’*éŒ…¡¥h™¶Ä!HÊY¥ò!Ò¤KiC7ñILÉZÄ±#˜HÇb¾7Ãy3³«q4Š“7g!	“c‡’Ç$§t_±ÕŒä-Eh+¦¥ñŠ3{Î÷
µšrø$Ø¥àÚ€ÕØŞ½UqSH¢¶š¯s ÉkèğJVØ¡—\;¥‡;óD!†ß³º5P´s
¥Ğ9xNÖò%ï$‡3÷y†BÃƒ®©aeÄş.=İø°á*<ùˆìšë1ãLg$ô¡£‰®0h°œÉ(M|Èø@š¶94óJšLhST8³!I¼¾›N¬Ú˜N²@íş(xuzó?¯´úl×NãMêÓ²Ó-H,ùK³BêÕrJf!K¾®Ü0O“’E^éÇ{¶75VÄË:êÆÍiXœû?>¦^R7rª³—(_ö!Š#Ÿ>R7{³Äd„êæRš½Â÷²°%|)ôu‰ß%X®!‡»³Cc·³Ó½X™õ[ tC	 ³Üš¥‹ô¾e­Û`w*·zåıí¯B¢ìaóÜßş„¡·…â`Üqpy–”¾Í…ğÙ™R¨ë¡ÏÿáÙÑ¯Z›i
```

## .git/objects/6d/6408a5bb8411448f2ee0f3e05078f8e13ff684

> Fichier binaire non inclus (1673 octets)

## .git/objects/87/63c1d9aa8a5da1da10bdcac57f1fc5ebb09491

> Fichier binaire non inclus (55 octets)

## .git/objects/9d/c71497ae68a4e83a6d720c3ad2fe6abd38d86f

> Fichier binaire non inclus (305 octets)

## .git/objects/9d/cf46457a2239e54fbe4a9f00f2826996522e04

> Fichier binaire non inclus (168 octets)

## .git/objects/a2/91f3e344cf5c4aef3053ba2cb40e91b5fc39b3

> Fichier binaire non inclus (85 octets)

## .git/objects/ac/d9d78e43721d13834d5d0f2c10a60001a496b6

> Fichier binaire non inclus (215 octets)

## .git/objects/b7/6e4ffecf8364809c61229d2644b578f6a94157

> Fichier binaire non inclus (88 octets)

## .git/objects/be/0c220c16ac5a2391490b3f9618e5f9454345aa

> Fichier binaire non inclus (527 octets)

## .git/objects/c2/7fb11e9c9c3fe264ad1fdf894bbb8214c1d6a7

```text
x5ÁJ1†=ç)~ÖCÌ®¤àIH+,ÒF‚"R<¤ÙX"1’-ÔGò9|1“®^†™ÿŸùffçi‡ùÍâì^É5öv~Íµ.XÆ^¤zXö
‘1õ¼AˆŸp!Ú{ÜíS4]ü#ÿU1ÑHÜ“lbì‚¢³8¸0F{-’6eÅ%òñu2»üL;õµu–s!•Z‰§^nĞ‹¾JÅ9ş|÷^Ø¾à)gg2’±˜Õ7ºÙ„>å.öØ–#W¬—Ø6Un®Ğü;ÍÛ/KT6
```

## .git/objects/d7/0c9ad928febd9287bcb9079f5857db1af09c06

> Fichier binaire non inclus (143 octets)

## .git/objects/e2/b928d9112e1dc2bf98d9006e9a0bf2fc8cbf93

```text
xeRmkÔ@ös~ÅP´Ùp¹mŠV¤ñŠÈ)ÒÛoG¡›Ü$¬ÆİswÃ!GşÃ?Ö™Í½YçÃ2;óÌ<óVu¶‚w¯ß¾¨­ñZ·ªaõÚ¡H?°áŒŸéwŸfe2ÂVÎûÕª%ºÿÑÑ9í¢÷(DÕ?T‹sl´ÑA[CGi$ãï~›Z¤^]‡h¶İ2)Eq5~û<¥X=·d’#Ë$izSGV×›Ïº²"ƒM$u :@J~^Få=iET'“’ÑT+LÃM³°oÊòÌU@iìZĞĞØË²çÒœ¢Äó¦`¡Ï¿LŒ_(by ì¢Ê‰vBŸC2è„†WT+	Ìf”NO)ïÇe3–…K´ÒôZñ¸¸¾™z€oÔ÷Ëu?À”=€è<}îñçÊ_’Mš!sì0;cÎÁ?µÇÃ¾ÙQ’!ÙMİi4<ƒkØ/T~TÎit"mh-ÓÖ^^ÅÅyš{®.)J«ÎKÒi¨_ŒÇº§ãÌ¶ç¸kèäÚ.‘nêïŸP!P$Â²¯[£:)å	ÇÄ:äZépoïx]b3ä Ğ¹fWÛ»àÑFËnfô±.šÆ`ç‘oe¼¦2²ò	Pü0
```

## .git/objects/ef/327c55831b06abfc2b987608b12db85728ffbf

```text
xURİJÃ0õºOq(ÈRìJ+›Ha^‚W:Ü@D$$]ÚÖ´$2Æ^È×ğÅL¶n­ßEÎ9ßø¦á¸Ëî¯dİ6Ú¢Òmto+k”º©avª -¿EÇ<Õ­İı§¨÷<óLk)ôÂny¬D	½U´”¼!QÀYÙhA*h¦*A²YÖ‘^ÀbpÌ:øèà1c™¶ÔWå%ş ÑEà£Ê>æ$u6:Ìc0Ü€_|='Kç~Ìûaæ²ƒ©•ƒÊö¢³µZ*KÊğsş±|~}ùÂ›ëgïš:`Œ½<@Z¡{/EİšûAÅí›É“Iy0a.?¥ŠÕ‚RŸ>¤´fRQ²K»>.(‘Êˆb«-ÖL)±!#?ÜqÕäÓ4f£Ì ãúÊ[‡›Ü`9¤Óô#<uÎwvİ(´ú÷ÇÆ
ÌZ¡¬Àj#+Å6I’¸’ÏÃğ¡“w&í²Yø-‘ã!Q¯èÀ¿y¶S
```

## .git/objects/f5/c1c09205f17763c8289c63c6d35fd4df64c6a8

> Fichier binaire non inclus (126 octets)

## .git/refs/heads/master

```text
9dcf46457a2239e54fbe4a9f00f2826996522e04

```

## .git/refs/remotes/origin/master

```text
9dcf46457a2239e54fbe4a9f00f2826996522e04

```

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
  fibo-go:
    build:
      context: .
      dockerfile: go/Dockerfile
    container_name: fibo-go
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

  fibo-python:
    build:
      context: .
      dockerfile: python/Dockerfile
    depends_on: [fibo-go]
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

  fibo-node:
    build:
      context: .
      dockerfile: node/Dockerfile
    depends_on: [fibo-go]
    deploy:
      resources:
        limits:
          cpus: '5.0'
          memory: 5G
    networks: [fibo-net]

networks:
  fibo-net:
    driver: bridge
```

## go/Dockerfile

```text
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache protoc protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /src

# 1. PrÃ©paration des modules
COPY go/go.mod ./
RUN go mod download

# 2. GÃ©nÃ©ration gRPC
# On copie le proto
COPY proto/sync.proto ./proto/

# On gÃ©nÃ¨re le code en disant Ã  protoc que la racine du module est '.'
# Cela va crÃ©er automatiquement le dossier /src/pb/
RUN mkdir -p pb
RUN protoc --proto_path=./proto \
    --go_out=. --go_opt=module=fibonacci \
    --go-grpc_out=. --go-grpc_opt=module=fibonacci \
    ./proto/sync.proto

# 3. Build du binaire
COPY go/main.go .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/fibo main.go

# --- Image finale ---
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fibo .
EXPOSE 50051
ENTRYPOINT ["./fibo"]
```

## go/go.mod

```text
module fibonacci

go 1.25.1

require google.golang.org/grpc v1.77.0

require (
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

```

## go/go.sum

```text
github.com/go-logr/logr v1.4.3 h1:CjnDlHq8ikf6E492q6eKboGOC0T8CDaOvkHCIg8idEI=
github.com/go-logr/logr v1.4.3/go.mod h1:9T104GzyrTigFIr8wt5mBrctHMim0Nb2HLGrmQ40KvY=
github.com/go-logr/stdr v1.2.2 h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=
github.com/go-logr/stdr v1.2.2/go.mod h1:mMo/vtBO5dYbehREoey6XUKy/eSumjCCveDpRre4VKE=
github.com/golang/protobuf v1.5.4 h1:i7eJL8qZTpSEXOPTxNKhASYpMn+8e5Q6AdndVa1dWek=
github.com/golang/protobuf v1.5.4/go.mod h1:lnTiLA8Wa4RWRcIUkrtSVa5nRhsEGBg48fD6rSs7xps=
github.com/google/go-cmp v0.7.0 h1:wk8382ETsv4JYUZwIsn6YpYiWiBsYLSJiTsyBybVuN8=
github.com/google/go-cmp v0.7.0/go.mod h1:pXiqmnSA92OHEEa9HXL2W4E7lf9JzCmGVUdgjX3N/iU=
github.com/google/uuid v1.6.0 h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
go.opentelemetry.io/auto/sdk v1.2.1 h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=
go.opentelemetry.io/auto/sdk v1.2.1/go.mod h1:KRTj+aOaElaLi+wW1kO/DZRXwkF4C5xPbEe3ZiIhN7Y=
go.opentelemetry.io/otel v1.38.0 h1:RkfdswUDRimDg0m2Az18RKOsnI8UDzppJAtj01/Ymk8=
go.opentelemetry.io/otel v1.38.0/go.mod h1:zcmtmQ1+YmQM9wrNsTGV/q/uyusom3P8RxwExxkZhjM=
go.opentelemetry.io/otel/metric v1.38.0 h1:Kl6lzIYGAh5M159u9NgiRkmoMKjvbsKtYRwgfrA6WpA=
go.opentelemetry.io/otel/metric v1.38.0/go.mod h1:kB5n/QoRM8YwmUahxvI3bO34eVtQf2i4utNVLr9gEmI=
go.opentelemetry.io/otel/sdk v1.38.0 h1:l48sr5YbNf2hpCUj/FoGhW9yDkl+Ma+LrVl8qaM5b+E=
go.opentelemetry.io/otel/sdk v1.38.0/go.mod h1:ghmNdGlVemJI3+ZB5iDEuk4bWA3GkTpW+DOoZMYBVVg=
go.opentelemetry.io/otel/sdk/metric v1.38.0 h1:aSH66iL0aZqo//xXzQLYozmWrXxyFkBJ6qT5wthqPoM=
go.opentelemetry.io/otel/sdk/metric v1.38.0/go.mod h1:dg9PBnW9XdQ1Hd6ZnRz689CbtrUp0wMMs9iPcgT9EZA=
go.opentelemetry.io/otel/trace v1.38.0 h1:Fxk5bKrDZJUH+AMyyIXGcFAPah0oRcT+LuNtJrmcNLE=
go.opentelemetry.io/otel/trace v1.38.0/go.mod h1:j1P9ivuFsTceSWe1oY+EeW3sc+Pp42sO++GHkg4wwhs=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 h1:6/3JGEh1C88g7m+qzzTbl3A0FtsLguXieqofVLU/JAo=
golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82/go.mod h1:Q9BGdFy1y4nkUwiLvT5qtyhAnEHgnQ/zd8PfU6nc210=
golang.org/x/sys v0.37.0 h1:fdNQudmxPjkdUTPnLn5mdQv7Zwvbvpaxqs831goi9kQ=
golang.org/x/sys v0.37.0/go.mod h1:OgkHotnGiDImocRcuBABYBEXf8A9a87e/uXjp9XT3ks=
golang.org/x/text v0.30.0 h1:yznKA/E9zq54KzlzBEAWn1NXSQ8DIp/NYMy88xJjl4k=
golang.org/x/text v0.30.0/go.mod h1:yDdHFIX9t+tORqspjENWgzaCVXgk0yYnYuSZ8UzzBVM=
gonum.org/v1/gonum v0.16.0 h1:5+ul4Swaf3ESvrOnidPp4GZbzf0mxVQpDCYUQE7OJfk=
gonum.org/v1/gonum v0.16.0/go.mod h1:fef3am4MQ93R2HHpKnLk4/Tbh/s0+wqD5nfa6Pnwy4E=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 h1:M1rk8KBnUsBDg1oPGHNCxG4vc1f49epmTO7xscSajMk=
google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8/go.mod h1:7i2o+ce6H/6BluujYR+kqX3GKH+dChPTQU19wjRPiGk=
google.golang.org/grpc v1.77.0 h1:wVVY6/8cGA6vvffn+wWK5ToddbgdU3d8MNENr4evgXM=
google.golang.org/grpc v1.77.0/go.mod h1:z0BY1iVj0q8E1uSQCjL9cppRj+gnZjzDnzV0dHhrNig=
google.golang.org/protobuf v1.36.10 h1:AYd7cD/uASjIL6Q9LiTjz8JLcrh/88q5UObnmY3aOOE=
google.golang.org/protobuf v1.36.10/go.mod h1:HTf+CrKn2C3g5S8VImy6tdcUvCska2kB7j23XfzDpco=

```

## go/main.go

```go
package main

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	// L'import doit correspondre au nom dans go.mod + le dossier de destination
	pb "fibonacci/pb"
)

type server struct {
	pb.UnimplementedBarrierServer
	mu         sync.Mutex
	cond       *sync.Cond
	readyCount int
}

// WaitToStart bloque les clients jusqu'Ã  ce que les 2 (Node + Python) soient lÃ 
func (s *server) WaitToStart(ctx context.Context, in *pb.Empty) (*pb.StartSignal, error) {
	s.mu.Lock()
	s.readyCount++
	fmt.Printf("[gRPC] Client connectÃ©. Total : %d/2\n", s.readyCount)
	if s.readyCount >= 2 {
		s.cond.Broadcast()
	} else {
		s.cond.Wait()
	}
	s.mu.Unlock()
	return &pb.StartSignal{Message: "Signal de dÃ©part reÃ§u !"}, nil
}

func runFibo() {
	for run := 1; run <= 10; run++ {
		fmt.Printf("[GO] --- DÃ©marrage Run %d/10 ---\n", run)
		a, b := big.NewInt(0), big.NewInt(1)
		start := time.Now()

		for i := 0; i <= 400000; i++ {
			a.Add(a, b)
			a, b = b, a

			if i%10000 == 0 && i > 0 {
				fmt.Printf("[GO] Run %d | %d itÃ©rations | Temps Ã©coulÃ©: %v\n", run, i, time.Since(start))
			}
		}
		fmt.Printf("[GO] Run %d terminÃ© en %v\n", run, time.Since(start))
	}
}

func main() {
	// 1. Initialisation du serveur gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := &server{}
	s.cond = sync.NewCond(&s.mu)
	grpcServer := grpc.NewServer()
	pb.RegisterBarrierServer(grpcServer, s)

	// 2. Lancement du serveur en arriÃ¨re-plan
	go func() {
		fmt.Println("[GO] Serveur gRPC Ã  l'Ã©coute sur :50051")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// 3. Attente de la barriÃ¨re de synchronisation pour le worker GO lui-mÃªme
	fmt.Println("[GO] En attente de Node.js et Python...")
	s.mu.Lock()
	for s.readyCount < 2 {
		s.cond.Wait()
	}
	s.mu.Unlock()

	// 4. Lancement du calcul
	fmt.Println("[GO] Signal reÃ§u. Lancement des calculs !")
	runFibo()

	// On laisse un peu de temps pour que les autres finissent avant de couper le container
	time.Sleep(30 * time.Second)
}

```

## main.go

```go
package main

import (
	"fmt"
	"math/big"
)

// FibGenerator produit la suite de Fibonacci via un channel.
// S'arrÃªte si la RAM estimÃ©e des deux derniers nombres dÃ©passe limitBytes.
func FibGenerator(limitBytes uint64) <-chan *big.Int {
	ch := make(chan *big.Int)

	go func() {
		defer close(ch)
		// Initialisation : F(0)=0, F(1)=1
		a := big.NewInt(0)
		b := big.NewInt(1)

		for {
			// Calcul de la taille approximative en RAM des deux termes
			// big.Int stocke les donnÃ©es dans un slice de 'Word' (uint sur 64 bits)
			// On compte environ 8 octets par Word + le overhead de la structure.
			sizeA := uint64(len(a.Bits())) * 8
			sizeB := uint64(len(b.Bits())) * 8

			if sizeA+sizeB > limitBytes {
				fmt.Printf("\n[Limite de %d Go atteinte]\n", limitBytes/1024/1024/1024)
				return
			}

			// On envoie une copie pour Ã©viter les effets de bord si l'appelant modifie la valeur
			val := new(big.Int).Set(a)
			ch <- val

			// Fibonacci : a, b = b, a+b
			// On utilise Add pour additionner b Ã  a, puis on swap.
			a.Add(a, b)
			a, b = b, a
		}
	}()

	return ch
}

func main() {
	const maxRAM = 5 * 1024 * 1024 * 1024 // 5 Go
	gen := FibGenerator(maxRAM)

	count := 0
	for f := range gen {
		count++
		// Pour l'exemple, on affiche tous les 100 000 termes
		// car l'affichage console est trÃ¨s lent pour de gros nombres.
		if count%100000 == 0 {
			fmt.Printf("Terme nÂ°%d calculÃ© (Taille actuelle : ~%d MB)\n", count, len(f.Bits())*8/1024/1024)
		}
	}
}

```

## node/Dockerfile

```text
FROM node:20-alpine

WORKDIR /app

RUN npm install @grpc/grpc-js @grpc/proto-loader

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# --- CORRECTION ICI ---
# On spÃ©cifie le dossier source 'node/'
COPY node/index.js .

CMD ["node", "index.js"]
```

## node/index.js

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
console.log("Node prÃªt, en attente du signal...");
client.waitToStart({}, (err) => {
    if (err) console.error(err);
    else runFibo();
});
```

## project_export.log

```text
[2025-12-25 15:13:06] Source  : .
[2025-12-25 15:13:06] Sortie  : project_export.md
[2025-12-25 15:13:06] Fichiers trouvÃ©s (avant filtre): 59
[2025-12-25 15:13:06] Fichiers Ã  concatÃ©ner (aprÃ¨s filtre): 58 (exclus auto:1 dir:0 file:0)
[2025-12-25 15:13:06] ConcatÃ¨ne [1] .git/COMMIT_EDITMSG (size=64)
[2025-12-25 15:13:06] ConcatÃ¨ne [2] .git/HEAD (size=23)
[2025-12-25 15:13:06] ConcatÃ¨ne [3] .git/config (size=258)
[2025-12-25 15:13:06] ConcatÃ¨ne [4] .git/description (size=73)
[2025-12-25 15:13:06] ConcatÃ¨ne [5] .git/hooks/applypatch-msg.sample (size=478)
[2025-12-25 15:13:06] ConcatÃ¨ne [6] .git/hooks/commit-msg.sample (size=896)
[2025-12-25 15:13:06] ConcatÃ¨ne [7] .git/hooks/fsmonitor-watchman.sample (size=4726)
[2025-12-25 15:13:06] ConcatÃ¨ne [8] .git/hooks/post-update.sample (size=189)
[2025-12-25 15:13:06] ConcatÃ¨ne [9] .git/hooks/pre-applypatch.sample (size=424)
[2025-12-25 15:13:06] ConcatÃ¨ne [10] .git/hooks/pre-commit.sample (size=1643)
[2025-12-25 15:13:06] ConcatÃ¨ne [11] .git/hooks/pre-merge-commit.sample (size=416)
[2025-12-25 15:13:06] ConcatÃ¨ne [12] .git/hooks/pre-push.sample (size=1374)
[2025-12-25 15:13:06] ConcatÃ¨ne [13] .git/hooks/pre-rebase.sample (size=4898)
[2025-12-25 15:13:06] ConcatÃ¨ne [14] .git/hooks/pre-receive.sample (size=544)
[2025-12-25 15:13:06] ConcatÃ¨ne [15] .git/hooks/prepare-commit-msg.sample (size=1492)
[2025-12-25 15:13:06] ConcatÃ¨ne [16] .git/hooks/push-to-checkout.sample (size=2783)
[2025-12-25 15:13:06] ConcatÃ¨ne [17] .git/hooks/sendemail-validate.sample (size=2308)
[2025-12-25 15:13:06] ConcatÃ¨ne [18] .git/hooks/update.sample (size=3650)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/index â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [20] .git/info/exclude (size=240)
[2025-12-25 15:13:06] ConcatÃ¨ne [21] .git/logs/HEAD (size=211)
[2025-12-25 15:13:06] ConcatÃ¨ne [22] .git/logs/refs/heads/master (size=211)
[2025-12-25 15:13:06] ConcatÃ¨ne [23] .git/logs/refs/remotes/origin/master (size=144)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/12/a122ea7678e77dcb0ffc86aeca9c7d69bd816d â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [25] .git/objects/20/708f6a2ef877be2ccc2101160dbce8791b4604 (size=273)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/37/ae47ae280e513c8aa2129bd8132aefba977851 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/5a/50d4c26668c2b5b465756529b5868e7f109e11 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [28] .git/objects/5d/ab255ae1904e7a27735e3c9641a108451695d3 (size=206)
[2025-12-25 15:13:06] ConcatÃ¨ne [29] .git/objects/5e/1af68ac259b3f3f34291e2803d602d13ce9fd5 (size=206)
[2025-12-25 15:13:06] ConcatÃ¨ne [30] .git/objects/5f/7131495b2d34d83cb08d08a31f739f2fa24949 (size=475)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/6d/6408a5bb8411448f2ee0f3e05078f8e13ff684 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/87/63c1d9aa8a5da1da10bdcac57f1fc5ebb09491 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/9d/c71497ae68a4e83a6d720c3ad2fe6abd38d86f â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/9d/cf46457a2239e54fbe4a9f00f2826996522e04 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/a2/91f3e344cf5c4aef3053ba2cb40e91b5fc39b3 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/ac/d9d78e43721d13834d5d0f2c10a60001a496b6 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/b7/6e4ffecf8364809c61229d2644b578f6a94157 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/be/0c220c16ac5a2391490b3f9618e5f9454345aa â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [39] .git/objects/c2/7fb11e9c9c3fe264ad1fdf894bbb8214c1d6a7 (size=216)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/d7/0c9ad928febd9287bcb9079f5857db1af09c06 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [41] .git/objects/e2/b928d9112e1dc2bf98d9006e9a0bf2fc8cbf93 (size=470)
[2025-12-25 15:13:06] ConcatÃ¨ne [42] .git/objects/ef/327c55831b06abfc2b987608b12db85728ffbf (size=365)
[2025-12-25 15:13:06] â„¹ï¸  Binaire : .git/objects/f5/c1c09205f17763c8289c63c6d35fd4df64c6a8 â€” rÃ©fÃ©rencÃ© mais non inclus
[2025-12-25 15:13:06] ConcatÃ¨ne [44] .git/refs/heads/master (size=41)
[2025-12-25 15:13:06] ConcatÃ¨ne [45] .git/refs/remotes/origin/master (size=41)
[2025-12-25 15:13:06] ConcatÃ¨ne [46] .gitignore (size=123)
[2025-12-25 15:13:06] ConcatÃ¨ne [47] compose.yml (size=697)
[2025-12-25 15:13:06] ConcatÃ¨ne [48] go/Dockerfile (size=909)
[2025-12-25 15:13:06] ConcatÃ¨ne [49] go/go.mod (size=365)
[2025-12-25 15:13:06] ConcatÃ¨ne [50] go/go.sum (size=3182)
[2025-12-25 15:13:06] ConcatÃ¨ne [51] go/main.go (size=2057)
[2025-12-25 15:13:06] ConcatÃ¨ne [52] main.go (size=1476)
[2025-12-25 15:13:06] ConcatÃ¨ne [53] node/Dockerfile (size=257)
[2025-12-25 15:13:06] ConcatÃ¨ne [54] node/index.js (size=836)

```

## proto/sync.proto

```text
syntax = "proto3";

package sync;

// Indique que le package fait partie du module 'fibonacci' dans le dossier 'pb'
option go_package = "fibonacci/pb";

service Barrier {
  rpc WaitToStart (Empty) returns (StartSignal);
}

message Empty {}
message StartSignal {
  string message = 1;
}
```

## python/Dockerfile

```text
FROM python:3.12-slim

WORKDIR /app

RUN pip install --no-cache-dir grpcio grpcio-tools

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# GÃ©nÃ©ration du code Python
RUN python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. sync.proto

# --- CORRECTION ICI ---
# On spÃ©cifie le dossier source 'python/'
COPY python/main.py .

CMD ["python", "main.py"]
```

## python/main.py

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
        print("Python prÃªt, en attente du signal...")
        stub.WaitToStart(Empty())
        run_fibo()
```

