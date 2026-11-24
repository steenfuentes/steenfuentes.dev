module.exports = async ({ github, context, core, previewFrontendUrl, previewBackendUrl }) => {
  const body = `🚀 **Preview deployment ready!**

          🌐 Frontend: ${previewFrontendUrl}
          🔧 Backend API: ${previewBackendUrl}

          This preview will be automatically deleted when the PR is closed.`;

          // Find existing comment
          const { data: comments } = await github.rest.issues.listComments({
            owner: context.repo.owner,
            repo: context.repo.repo,
            issue_number: context.issue.number,
          });

          const botComment = comments.find(comment =>
            comment.user.type === 'Bot' &&
            comment.body.includes('Preview deployment ready!')
          );

          if (botComment) {
            await github.rest.issues.updateComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              comment_id: botComment.id,
              body: body
            });
          } else {
            await github.rest.issues.createComment({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              body: body
            });
          }
        }