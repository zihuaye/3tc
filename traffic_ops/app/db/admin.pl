#!/usr/bin/env perl
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
use strict;
use File::Basename;

STDERR->autoflush(1);
print STDERR "WARNING: this script is deprecated, please use the db/admin binary instead.\n\n";
my $args = join(" ", @ARGV);
my $new_admin = dirname(__FILE__) . "/admin";
exit(system("$new_admin " . $args));
