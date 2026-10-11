**CHECK FOR DUPLICATE ISSUES BEFORE OPENING THIS ONE**

Stale and contradictory wording in [[Fuzzing tutorial > Fix the double reverse error > Write the code > Step 2]](https://go.dev/doc/tutorial/fuzz#write-the-code-6)

- **Source:**
  - Tutorial: [go.dev/doc/tutorial/fuzz#write-the-code-6](https://go.dev/doc/tutorial/fuzz#write-the-code-6)
  - Repository file: [`_content/doc/tutorial/fuzz.md`](https://github.com/golang/website/blob/master/_content/doc/tutorial/fuzz.md)
    - Target section: [L606–L622](https://github.com/golang/website/blob/master/_content/doc/tutorial/fuzz.md#L606-L622)
  - Blame: [`_content/doc/tutorial/fuzz.md#L607`](https://github.com/golang/website/blame/bb3c125582905e8bfaab7fb1775d1a5afecc9d13/_content/doc/tutorial/fuzz.md#L607)
  - Historical context: Gerrit CL [381794](https://go-review.googlesource.com/c/website/+/381794) / commit [`feaf54d4`](https://github.com/golang/website/commit/feaf54d44c0c3939771de7c81c9f95bf94a50d8d) (*"_content/doc/tutorial: add additional descriptions to fuzz.md"* by Katie Hockman, Feb 2022).

---

### Issue Summary

In the **"Fix the double reverse error"** section of the Fuzzing tutorial, Step 2 instructs the learner to *"discard the extra error value"*, but the provided replacement code does the opposite: it assigns the error return values to variables and logs them.

In Go idiom, "discarding" an error specifically means assigning it to the blank identifier `_` (e.g., `rev, _ := Reverse(input)`). Instructing beginners to "discard" an error while providing code that captures and formats it creates conceptual confusion for learners.

---

### Historical Root Cause

Before Gerrit CL [381794](https://go-review.googlesource.com/c/website/+/381794), the tutorial code actually discarded the errors using `_`:

```go
// Code prior to CL 381794:
func main() {
    input := "The quick brown fox jumped over the lazy dog"
    rev, _ := Reverse(input)
    doubleRev, _ := Reverse(rev)
    fmt.Printf("original: %q\n", input)
    fmt.Printf("reversed: %q\n", rev)
    fmt.Printf("reversed again: %q\n", doubleRev)
}
```

At that time, the instruction *"modify the `main` function to discard the extra error value"* accurately described the code.

In CL 381794, the code snippet was updated to capture `revErr` and `doubleRevErr` and print them via `fmt.Printf`, but the introductory sentence was accidentally left unchanged.

---

### Comparison

#### Currently Documented in Tutorial

```markdown
1. Since the Reverse function now returns an error, modify the `main` function to
   discard the extra error value. Replace the existing `main` function with the
   following.

   ```
   func main() {
       input := "The quick brown fox jumped over the lazy dog"
       rev, revErr := Reverse(input)
       doubleRev, doubleRevErr := Reverse(rev)
       fmt.Printf("original: %q\n", input)
       fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
       fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
   }
   ```
```

#### Discrepancy

* **Narrative says:** Discard the extra error value.
* **Code actually does:** Captures `revErr` and `doubleRevErr`, then prints them (`err: %v`). Nothing is discarded.

---

### Suggested Fix

#### Approach A (Recommended: Update narrative to match code)

Update the introductory sentence to reflect what the code actually does:

```markdown
Since the Reverse function now returns an error, modify the `main` function to
capture and print the error values. Replace the existing `main` function with the
following.
```

*(Alternatively: `... modify the main function to handle and print the returned errors.`)*

#### Approach B (Revert code snippet to discard error)

If the intent was indeed to keep `main` simple by discarding the error:

```go
func main() {
    input := "The quick brown fox jumped over the lazy dog"
    rev, _ := Reverse(input)
    doubleRev, _ := Reverse(rev)
    fmt.Printf("original: %q\n", input)
    fmt.Printf("reversed: %q\n", rev)
    fmt.Printf("reversed again: %q\n", doubleRev)
}
```

> **Note:** Approach A is preferred because printing errors is idiomatic and reinforces good error handling practices for learners.