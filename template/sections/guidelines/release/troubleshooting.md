### Troubleshooting

#### Workflow Failed

1. Check workflow logs:

   ```bash
   gh run view <run-id> --log-failed
   ```

2. Common issues:
   - **Build errors**: Run `make check-build` locally
   - **CGO issues**: Ensure code is CGO-free compatible
   - **Tag already exists**: Delete and recreate if needed

#### Delete a Tag

If you need to delete a tag (e.g., wrong version):

```bash
# Delete local tag
git tag -d v6.1.0

# Delete remote tag
git push origin --delete v6.1.0
```

**Note**: If a release was created, delete it first on GitHub.

#### Re-run Failed Release

If the workflow failed after tag was pushed:

1. Fix the issue in code
2. Delete the tag (local and remote)
3. Create a new commit with the fix
4. Create and push the tag again
