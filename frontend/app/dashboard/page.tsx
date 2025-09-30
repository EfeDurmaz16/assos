'use client'

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getDashboardData, createVideo } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { 
  Play, 
  Plus, 
  Bot, 
  TrendingUp, 
  Users, 
  DollarSign, 
  Video, 
  Clock,
  Eye,
  Heart,
  MessageCircle,
  Share,
  Settings,
  BarChart3,
  AlertCircle
} from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'

function StatCard({ title, value, icon: Icon, isLoading }: any) {
  return (
    <Card>
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium text-gray-600">{title}</p>
            {isLoading ? <Skeleton className="h-7 w-24 mt-1" /> : <p className="text-2xl font-bold">{value}</p>}
          </div>
          <Icon className="w-8 h-8 text-gray-400" />
        </div>
      </CardContent>
    </Card>
  )
}

export default function DashboardPage() {
  const [isCreating, setIsCreating] = useState(false);

  const { data, error, isLoading } = useQuery({
    queryKey: ['dashboardData'],
    queryFn: getDashboardData,
    refetchInterval: 15000, // Refetch every 15 seconds
  });

  const handleCreateVideo = async () => {
    setIsCreating(true);
    try {
      // This is a placeholder for a real creation flow which would likely involve a modal/form
      const videoData = {
        channel_id: "your-default-channel-id", // Replace with actual channel ID from state/context
        title: "New AI-Generated Video",
        description: "This video was created automatically."
      };
      const newVideo = await createVideo(videoData);
      // In a real app, you'd probably show a success toast and refetch the videos list
      console.log("Created video:", newVideo);
    } catch (err) {
      console.error("Failed to create video:", err);
      // Show an error toast
    } finally {
      setIsCreating(false);
    }
  };

  const stats = data?.stats || {};
  const recentVideos = data?.recent_videos || [];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b">
        <div className="px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 gradient-youtube rounded-lg flex items-center justify-center">
                  <Play className="w-4 h-4 text-white" />
                </div>
                <span className="text-2xl font-bold text-gradient">ASSOS</span>
              </div>
              
              <nav className="hidden md:flex items-center space-x-6">
                <span className="text-gray-900 font-medium">Dashboard</span>
                <a href="#" className="text-gray-600 hover:text-gray-900">Channels</a>
                <a href="#" className="text-gray-600 hover:text-gray-900">Videos</a>
                <a href="#" className="text-gray-600 hover:text-gray-900">Analytics</a>
                <a href="#" className="text-gray-600 hover:text-gray-900">AI Agents</a>
              </nav>
            </div>

            <div className="flex items-center space-x-4">
              <Button 
                className="gradient-youtube text-white"
                onClick={handleCreateVideo}
                disabled={isCreating}
              >
                <Plus className="w-4 h-4 mr-2" />
                {isCreating ? 'Creating...' : 'Create Video'}
              </Button>
              <Button variant="outline" size="icon">
                <Settings className="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>
      </header>

      <div className="p-6">
        {/* Stats Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6 mb-8">
          <StatCard title="Total Videos" value={stats.total_videos} icon={Video} isLoading={isLoading} />
          <StatCard title="Processing" value={stats.processing_videos} icon={Clock} isLoading={isLoading} />
          <StatCard title="Published" value={stats.published_videos} icon={TrendingUp} isLoading={isLoading} />
          <StatCard title="Total Views" value={stats.total_views?.toLocaleString()} icon={Eye} isLoading={isLoading} />
          <StatCard title="Revenue" value={`$${stats.total_revenue?.toFixed(2)}`} icon={DollarSign} isLoading={isLoading} />
          <StatCard title="Avg. CTR" value={`${stats.avg_ctr?.toFixed(1)}%`} icon={BarChart3} isLoading={isLoading} />
        </div>

        {error && (
          <Card className="mb-6 bg-red-50 border-red-200">
            <CardContent className="p-4 flex items-center">
              <AlertCircle className="w-5 h-5 text-red-600 mr-3" />
              <p className="text-sm text-red-700">
                Failed to load dashboard data. {(error as Error).message}
              </p>
            </CardContent>
          </Card>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Recent Videos */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Recent Videos</CardTitle>
                <CardDescription>
                  Your latest video content and processing status
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {isLoading ? (
                    Array.from({ length: 3 }).map((_, i) => (
                      <div key={i} className="flex items-center space-x-4 p-4">
                        <Skeleton className="w-16 h-10 rounded" />
                        <div className="flex-1 space-y-2">
                          <Skeleton className="h-4 w-3/4" />
                          <Skeleton className="h-3 w-1/2" />
                        </div>
                        <Skeleton className="h-6 w-20 rounded-full" />
                      </div>
                    ))
                  ) : (
                    recentVideos.map((video: any) => (
                      <div key={video.id} className="flex items-center space-x-4 p-4 border rounded-lg">
                        <div className="w-16 h-10 bg-gray-200 rounded flex items-center justify-center">
                          {video.status === 'published' ? (
                            <Play className="w-4 h-4 text-gray-600" />
                          ) : (
                            <Clock className="w-4 h-4 text-ai-purple animate-pulse" />
                          )}
                        </div>
                        
                        <div className="flex-1">
                          <h3 className="font-medium text-sm mb-1">{video.title}</h3>

                          {video.status === 'published' ? (
                            <div className="flex items-center space-x-4 text-xs text-gray-600">
                              <span className="flex items-center">
                                <Eye className="w-3 h-3 mr-1" />
                                {video.views?.toLocaleString() || 0}
                              </span>
                              <span className="flex items-center">
                                <Heart className="w-3 h-3 mr-1" />
                                {video.likes?.toLocaleString() || 0}
                              </span>
                              <span className="flex items-center">
                                <MessageCircle className="w-3 h-3 mr-1" />
                                {video.comments || 0}
                              </span>
                            </div>
                          ) : (
                            <div className="flex items-center space-x-2">
                              <div className="flex-1 bg-gray-200 rounded-full h-2">
                                <div
                                  className="bg-ai-purple h-2 rounded-full transition-all"
                                  style={{ width: `${video.progress || 0}%` }}
                                ></div>
                              </div>
                              <span className="text-xs text-gray-600">{video.progress || 0}%</span>
                            </div>
                          )}
                        </div>

                        <div className="text-right">
                          <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize ${
                            video.status === 'published'
                              ? 'bg-green-100 text-green-800'
                              : 'bg-blue-100 text-blue-800'
                          }`}>
                            {video.status}
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* AI Agents Status */}
          <div>
            <Card>
              <CardHeader>
                <CardTitle>AI Agents</CardTitle>
                <CardDescription>
                  Status of your AI automation agents
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <Bot className="w-5 h-5 text-ai-purple" />
                      <div>
                        <p className="font-medium text-sm">Manus Orchestrator</p>
                        <p className="text-xs text-gray-600">Strategic planning</p>
                      </div>
                    </div>
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  </div>

                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <TrendingUp className="w-5 h-5 text-automation-green" />
                      <div>
                        <p className="font-medium text-sm">Trend Predictor</p>
                        <p className="text-xs text-gray-600">Market analysis</p>
                      </div>
                    </div>
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  </div>

                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <BarChart3 className="w-5 h-5 text-youtube" />
                      <div>
                        <p className="font-medium text-sm">Content Strategist</p>
                        <p className="text-xs text-gray-600">Script generation</p>
                      </div>
                    </div>
                    <div className="w-2 h-2 bg-yellow-500 rounded-full pulse-ai"></div>
                  </div>

                  <div className="flex items-center justify-between p-3 border rounded-lg">
                    <div className="flex items-center space-x-3">
                      <Users className="w-5 h-5 text-ai-purple" />
                      <div>
                        <p className="font-medium text-sm">Research Agent</p>
                        <p className="text-xs text-gray-600">Content research</p>
                      </div>
                    </div>
                    <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  </div>
                </div>

                <Button variant="outline" className="w-full mt-4">
                  <Settings className="w-4 h-4 mr-2" />
                  Manage Agents
                </Button>
              </CardContent>
            </Card>

            {/* Quick Actions */}
            <Card className="mt-6">
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button variant="outline" className="w-full justify-start">
                  <Plus className="w-4 h-4 mr-2" />
                  Create New Video
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <BarChart3 className="w-4 h-4 mr-2" />
                  View Analytics
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <TrendingUp className="w-4 h-4 mr-2" />
                  Trend Analysis
                </Button>
                <Button variant="outline" className="w-full justify-start">
                  <Settings className="w-4 h-4 mr-2" />
                  Channel Settings
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}