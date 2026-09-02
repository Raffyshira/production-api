import { API_URL } from "@/lib/api";
import { useState } from "react";
import { useCookies } from "react-cookie";
import { Button } from "./components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "./components/ui/drawer";
import { Input } from "./components/ui/input";

import { useMediaQuery } from "@/hooks/use-media-query";

export const CreatePostForm: React.FC<{ onFetchPosts: () => void }> = ({
  onFetchPosts,
}) => {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [open, setOpen] = useState(false);

  const isDesktop = useMediaQuery("md");

  const [cookies] = useCookies(["at"]);
  const at = cookies.at;

  const handleSubmit = async () => {
    if (!title.trim() || !content.trim()) return;

    await fetch(`${API_URL}/posts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${at}`,
      },
      body: JSON.stringify({
        title,
        content,
      }),
    });

    setTitle("");
    setContent("");
    setOpen(false);
    onFetchPosts();
  };

  return (
    <Drawer
      open={open}
      onOpenChange={setOpen}
      showSwipeHandle
      swipeDirection={isDesktop ? "right" : "down"}
    >
      <DrawerTrigger
        render={
          <Button
            variant="outline"
            className="w-full justify-start text-muted-foreground"
          />
        }
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="mr-2"
        >
          <path d="M12 5v14" />
          <path d="M5 12h14" />
        </svg>
        What's on your mind?
      </DrawerTrigger>

      <DrawerContent>
        <div className="mx-auto w-full max-w-lg">
          <DrawerHeader>
            <DrawerTitle>Create a new post</DrawerTitle>
            <DrawerDescription>
              Share your thoughts with the community.
            </DrawerDescription>
          </DrawerHeader>

          <div className="flex flex-col gap-3 p-4">
            <Input
              placeholder="Post title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <textarea
              placeholder="What's on your mind..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={4}
              className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 resize-none"
            />
          </div>

          <DrawerFooter>
            <Button
              onClick={handleSubmit}
              disabled={!title.trim() || !content.trim()}
            >
              Share
            </Button>
            <DrawerClose render={<Button variant="outline" />}>
              Cancel
            </DrawerClose>
          </DrawerFooter>
        </div>
      </DrawerContent>
    </Drawer>
  );
};
