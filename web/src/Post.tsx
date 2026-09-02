interface PostComment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  created_at: string;
  user?: {
    id: number;
    username: string;
  };
}

export interface FeedPost {
  id: number;
  user_id: number;
  comments_count: number;
  content: string;
  created_at: string;
  tags: string[];
  title?: string;
  comments?: PostComment[];
  user?: {
    username: string;
  };
}

interface PostProps {
  post: FeedPost;
  onClick: () => void;
  onUnfollow?: (userId: number) => void;
}

export const Post: React.FC<PostProps> = ({ post, onClick, onUnfollow }) => {
  const date = new Date(post.created_at);
  const timeAgo = getTimeAgo(date);

  return (
    <article
      onClick={onClick}
      className="group cursor-pointer rounded-xl border bg-card p-5 transition-all duration-200 hover:border-foreground/20 hover:shadow-sm active:scale-[0.99]"
    >
      {/* Header: username + time + unfollow */}
      <div className="flex items-center justify-between text-sm text-muted-foreground">
        <div className="flex items-center gap-2">
          <div className="flex h-7 w-7 items-center justify-center rounded-full bg-secondary text-xs font-semibold text-secondary-foreground">
            {(post.user?.username ?? "?")[0].toUpperCase()}
          </div>
          <span className="font-medium text-foreground">
            {post.user?.username ?? "Unknown"}
          </span>
          <span>·</span>
          <span>{timeAgo}</span>
        </div>
        {onUnfollow && (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onUnfollow(post.user_id);
            }}
            className="text-xs text-muted-foreground hover:text-destructive hover:underline"
          >
            Unfollow
          </button>
        )}
      </div>

      {/* Title + content */}
      <div className="mt-3">
        <h2 className="text-base font-semibold leading-snug tracking-tight group-hover:text-foreground/80">
          {post.title || "Untitled"}
        </h2>
        <p className="mt-1 line-clamp-2 text-sm leading-relaxed text-muted-foreground">
          {post.content}
        </p>
      </div>

      {/* Tags */}
      {post.tags && post.tags.length > 0 && (
        <div className="mt-3 flex flex-wrap gap-1.5">
          {post.tags.map((tag) => (
            <span
              key={tag}
              className="rounded-full bg-secondary px-2 py-0.5 text-xs font-medium text-secondary-foreground"
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* Footer: comments count */}
      <div className="mt-3 flex items-center gap-1.5 text-xs text-muted-foreground transition-colors group-hover:text-foreground/60">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        </svg>
        <span>
          {post.comments_count}{" "}
          {post.comments_count === 1 ? "comment" : "comments"}
        </span>
      </div>
    </article>
  );
};

function getTimeAgo(date: Date): string {
  const now = new Date();
  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (seconds < 60) return "just now";

  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;

  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;

  const days = Math.floor(hours / 24);
  if (days < 7) return `${days}d ago`;

  return date.toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
  });
}
