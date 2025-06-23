"use client";

import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { getOwnerByName } from "@/api/owner.api";
import { EntityTypeEnum } from "@/types/entity";
import OrgView from "@/views/entities/OrgView";
import UserView from "@/views/entities/UserView";

interface OwnerProps {
  type: string;
  id: number;
}

export default function OwnerView() {
  const [ownerProps, setOwnerProps] = useState<OwnerProps | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const params = useParams();

  useEffect(() => {
    console.log("Owner params:", params); // 调试日志
    if (!params.owner) {
      console.error("No owner name provided");
      setError("无效的用户名");
      setLoading(false);
      return;
    }
    setLoading(true);
    getOwnerByName({ name: params.owner as string })
      .then((owner) => {
        console.log("Owner data:", owner);
        if (owner && owner.data) {
          setOwnerProps({
            type: owner.data.type,
            id: owner.data.id,
          });
        } else {
          setError("未找到用户");
        }
      })
      .catch((error) => {
        console.error("Error fetching owner:", error);
        setError("加载用户数据失败");
      })
      .finally(() => {
        setLoading(false);
      });
  }, [params.owner]);

  if (loading) {
    return (
      <div className="text-center p-8 text-slate-700 dark:text-slate-100">正在加载用户数据...</div>
    );
  }

  if (error) {
    return <div className="text-center p-8 text-red-500">{error}</div>;
  }

  return (
    <div className="container mx-auto p-4">
      {ownerProps ? (
        ownerProps.type === EntityTypeEnum.USER ? (
          <UserView id={ownerProps.id} />
        ) : (
          <OrgView id={ownerProps.id} />
        )
      ) : (
        <div className="text-center p-8 bg-yellow-100 dark:bg-yellow-900 rounded-lg">
          未找到用户数据
        </div>
      )}
    </div>
  );
}
