# boludo

**boludo** is a personal AI assistant that lives in your terminal. It respects
your privacy by running specialized [Large language models](https://en.wikipedia.org/wiki/Large_language_model)
directly on your ordinary computer.

`boludo` is a user-friendly interface for _[llama.cpp](https://github.com/ggerganov/llama.cpp)_.

## Usage

`boludo` uses subcommands to call locally hosted models with specific parameters.
`boludo.toml` is the config file with subcommands definitions (see:
[examples/boludo.toml](./examples/boludo.toml)). For example, if you defined
`someconfig`, you can do:

```sh
$ boludo someconfig "How are you?"
I am fine, thanks.
$ boludo someconfig <input.txt >output.txt
```

In the config file, you can change the default behaviour of the model by adjusting two
parameters:

- **creativity** modifies probabilities of the next _word_. For example:

   - `0.0 < creativity < 1.0` - increase predictability of result. The smaller the
      value, the more similar the result to the dataset used to train the model (useful
      for code generation or [plagiarizing _The New York Times_](https://nytco-assets.nytimes.com/2023/12/NYT_Complaint_Dec2023.pdf))
   - `creativity = 1.0` - preserves probabilities computed during model training
      (best for validating claims of the model author)
   - `1.0 < creativity` - increase probabilities of unlikely next
      _words_ (best for making the output more diverse, but higher values increase
      inconsistency).

   In popular LLM runners, this parameter is known as a **temperature**.

- **cutoff** discards statistically insignificant _words_ ([long tail](https://en.wikipedia.org/wiki/Long_tail))
   during the next _word_ prediction. For example, `cutoff = 0.03` means discarding
   _words_ with a probability smaller than 3% of the value of the most probable _word_.

   In popular LLM runners, this parameter is known as a **min-p**.

### Proofreading

_[Karen TheEditor](https://huggingface.co/FPHam/Karen_TheEditor_V2_STRICT_Mistral_7B)_
is one of the most practical models. It's specialized in correcting writing
(in American English). With a configuration similar to [examples/boludo.toml](./examples/boludo.toml),
you can do:

```sh
$ boludo proofreader "Trolling is a art"
Trolling is an art.
```

With longer texts, you may be interested to know which words were changed. By using the `tee` and `git` commands, you can monitor progress and track changes:

```sh
$ boludo proofreader <original.txt | tee corrected.txt; git -c color.diff.new="italic brightgreen reverse black bold" -c color.diff.old="strike brightred reverse bold" diff --no-index --color-words=. original.txt corrected.txt
```

### Coding

_[CodeNinja](https://huggingface.co/beowolx/CodeNinja-1.0-OpenChat-7B)_ is
a decent programmer assistant. With a configuration similar to [examples/boludo.toml](./examples/boludo.toml),
you can:

- find solutions for complex tasks:
   ```sh
   $ boludo coder "Show me the unix shell code without description for the task: Given a text file and an integer k, print the k most common words in the file (and the number of their occurrences) in decreasing frequency."

   cat text_file.txt | tr -s ' ' '\n' > words
   sort words | uniq -c > uniq_words
   awk '{if (length($0)<=k) print}' uniq_words > output.txt
   rm words uniq_words
   ```
   (not quite the [McIlroy's way](https://dl.acm.org/doi/10.1145/5948.315654))

- perform code completion:
   ```sh
   $ boludo coder "# returns anagram
   def anagram(str):"

      # sort the string
      sorted_str = ''.join(sorted(str.lower()))

      # return the result
      return sorted_str
   ```
   (it correctly recognized Python code)

- mentoring:
   ```sh
   $ boludo coder "Which 3 malpractices are most common among C# programmers?"
   1. Ignoring exceptions: Many C# programmers tend to ignore exceptions, which can lead to unhandled errors and unpredictable behavior in the application.

   2. Incorrect use of `string.Format`: Programmers often make mistakes when using `string.Format` for string manipulation, resulting in potential vulnerabilities or incorrect output.

   3. Lack of unit testing: Many C# developers neglect to write unit tests, which can lead to undiscovered bugs and poor code quality.

   To avoid these common mistakes, it's essential for C# programmers to follow best practices, utilize error-handling techniques, and practice thorough testing.
   ```
   (that's interesting - [format string attack](https://en.wikipedia.org/wiki/Uncontrolled_format_string)
   is mostly known in the C world, but it's still relevant in C#).

## Installation

You can manually build `boludo` with commands: `make && make build`.

## Development

Use `make` (GNU or BSD):

- `make` - install dependencies
- `make test` - runs test
- `make check` - static code analysis
- `make build` - compile binaries from latest commit
- `make dist` - compile binaries from latest commit for supported OSes
- `make clean` - removes compilation artifacts
- `make cli-release` - tag latest commit as a new release of CLI
- `make info` - print system info (useful for debugging).

### Versioning

The repo contains command-line utility which versions are tagged as `cli/vYYYY.0M.MICRO` (_[calendar versioning](https://calver.org/)_).

## TODO

Some ideas for further development:

- restrict allowed kernel calls (with _seccomp_ and _pledge_)
- embed `main` command from _llama.cpp_ using WASM (this will also make app more secure by sandboxing `main` command)
- implement native Go infering engine instead of _llama.cpp_ (the `gguf` format is well defined)
- show spinning wheel at stderr before generation.

## License

MIT ([in simple words](https://www.tldrlegal.com/license/mit-license))
