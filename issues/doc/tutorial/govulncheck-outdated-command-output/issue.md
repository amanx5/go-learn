**CHECK FOR DUPLICATE ISSUES BEFORE OPENING THIS ONE**

Outdated command output and narrative in [[Govulncheck tutorial > Install and run govulncheck > Steps 7 & 9]](https://go.dev/doc/tutorial/govulncheck#install-and-run-govulncheck)

- **Source:**
  - Tutorial: [go.dev/doc/tutorial/govulncheck](https://go.dev/doc/tutorial/govulncheck)
  - Repository file: [`_content/doc/tutorial/govulncheck.md`](https://github.com/golang/website/blob/237d41d9236bfaaabc9e42386a0816ae3b0e620d/_content/doc/tutorial/govulncheck.md)
    - Step 7 sample output: [L149–L188](https://github.com/golang/website/blob/237d41d9236bfaaabc9e42386a0816ae3b0e620d/_content/doc/tutorial/govulncheck.md#L149-L188)
    - Explanatory sections ("Interpreting the output", "Evaluate vulnerabilities"): [L190–L238](https://github.com/golang/website/blob/237d41d9236bfaaabc9e42386a0816ae3b0e620d/_content/doc/tutorial/govulncheck.md#L190-L238)
    - Step 9 sample output: [L265–L273](https://github.com/golang/website/blob/237d41d9236bfaaabc9e42386a0816ae3b0e620d/_content/doc/tutorial/govulncheck.md#L265-L273)
  - Historical context: Added in commit [`b3eae535`](https://github.com/golang/website/commit/b3eae535c5c6457317ee96c21685bf95f6ef259b) (April 2023) when `govulncheck` was pre-v1.0 (`@v0.0.0`) and experimental.

---

### Issue Summary

The command output demonstrated in Step 7 and Step 9 reflects pre-v1.0.0 `govulncheck` behavior from early 2023 and has diverged significantly from modern versions (v1.0.0+):

1. **Experimental Notice and Version Banner are Obsolete:**
   - The tutorial outputs:
     ```
     govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

     Using go1.20.3 and govulncheck@v0.0.0 with
     vulnerability data from https://vuln.go.dev (last modified 2023-04-18 21:32:26 +0000 UTC).

     Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
     ```
   - `govulncheck` officially reached general availability (GA) with **v1.0.0 in July 2023**. It no longer prints the experimental disclaimer or this header block during standard runs.

2. **`=== Informational ===` Section is Missing from Default Output:**
   - In Step 7, the tutorial shows an `=== Informational ===` section describing `GO-2022-1059` (an uncalled vulnerability in `golang.org/x/text@v0.3.5`).
   - Modern `govulncheck` hides uncalled vulnerabilities by default to keep noise low, providing only a summary line:
     ```
     This scan also found 1 vulnerability in packages you import and 14
     vulnerabilities in modules you require, but your code doesn't appear to call
     these vulnerabilities.
     Use '-show verbose' for more details.
     ```

3. **Narrative Disconnect in Subsequent Explanations:**
   - The subsequent subsections (*"Interpreting the output"* and *"Evaluate vulnerabilities"*) explicitly teach the reader how to interpret the `=== Informational ===` section and `GO-2022-1059`:
     > *"Another vulnerability, GO-2022-1059, exists in the golang.org/x/text module at v0.3.5. However, it is reported as 'Informational' because our code never (directly or indirectly) calls any of its vulnerable functions."*
     > *"We decided to deprioritize fixing the informational vulnerability, GO-2022-1059..."*
   - Learners following the tutorial today do not see `GO-2022-1059` or `=== Informational ===` in their terminal, making the narrative confusing and disconnected from the hands-on experience.
   - Additionally, when running modern `govulncheck -show verbose`, the section name `=== Informational ===` has been replaced by `=== Package Results ===` and `=== Module Results ===`.

4. **Step 9 (Post-Fix) Output is Also Outdated:**
   - Step 9 repeats the experimental banner and old summary format instead of the modern symbol result summary (`=== Symbol Results ===`, `No vulnerabilities found.`, `Your code is affected by 0 vulnerabilities.`).

---

### Comparison

#### Step 7 Output

- **Currently documented in tutorial:**
  ```
  govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

  Using go1.20.3 and govulncheck@v0.0.0 with
  vulnerability data from https://vuln.go.dev (last modified 2023-04-18 21:32:26 +0000 UTC).

  Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
  Your code is affected by 1 vulnerability from 1 module.

  Vulnerability #1: GO-2021-0113
    Due to improper index calculation, an incorrectly formatted
    language tag can cause Parse to panic via an out of bounds read.
    If Parse is used to process untrusted user inputs, this may be
    used as a vector for a denial of service attack.

    More info: https://pkg.go.dev/vuln/GO-2021-0113

    Module: golang.org/x/text
      Found in: golang.org/x/text@v0.3.5
      Fixed in: golang.org/x/text@v0.3.7

      Call stacks in your code:
        main.go:12:29: vuln.tutorial.main calls golang.org/x/text/language.Parse

  === Informational ===

  Found 1 vulnerability in packages that you import, but there are no call
  stacks leading to the use of this vulnerability. You may not need to
  take any action. See https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
  for details.

  Vulnerability #1: GO-2022-1059
    An attacker may cause a denial of service by crafting an
    Accept-Language header which ParseAcceptLanguage will take
    significant time to parse.
    More info: https://pkg.go.dev/vuln/GO-2022-1059
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.8
  ```

- **Actual modern output (`govulncheck ./...`):**
  ```
  === Symbol Results ===

  Vulnerability #1: GO-2021-0113
      Out-of-bounds read in golang.org/x/text/language
    More info: https://pkg.go.dev/vuln/GO-2021-0113
    Module: golang.org/x/text
      Found in: golang.org/x/text@v0.3.5
      Fixed in: golang.org/x/text@v0.3.7
      Example traces found:
        #1: main.go:14:29: vuln.tutorial.main calls language.Parse

  Your code is affected by 1 vulnerability from 1 module.
  This scan also found 1 vulnerability in packages you import and 14
  vulnerabilities in modules you require, but your code doesn't appear to call
  these vulnerabilities.
  Use '-show verbose' for more details.
  ```

---

#### Step 9 Output (After Upgrade to `v0.3.8`)

- **Currently documented in tutorial:**
  ```
  govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

  Using go1.20.3 and govulncheck@v0.0.0 with
  vulnerability data from https://vuln.go.dev (last modified 2023-04-06 19:19:26 +0000 UTC).

  Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
  No vulnerabilities found.
  ```

- **Actual modern output (`govulncheck ./...`):**
  ```
  === Symbol Results ===

  No vulnerabilities found.

  Your code is affected by 0 vulnerabilities.
  This scan also found 1 vulnerability in packages you import and 14
  vulnerabilities in modules you require, but your code doesn't appear to call
  these vulnerabilities.
  Use '-show verbose' for more details.
  ```

---

### Suggested Fix

1. Update the terminal output in Step 7 and Step 9 to match modern `govulncheck` formatting.
2. Address the discussion of `GO-2022-1059` using one of the following approaches:
   - **Approach A (Recommended):** Introduce `-show verbose` in the tutorial: explain that while `govulncheck` omits uncalled vulnerabilities by default, running `govulncheck -show verbose ./...` reveals uncalled package/module findings (including `GO-2022-1059` under `=== Package Results ===`), and tie that into the evaluation narrative.
   - **Approach B:** Streamline the narrative to focus on the summary message (*"This scan also found 1 vulnerability in packages you import..."*), removing specific textual references to `GO-2022-1059` and `=== Informational ===`.