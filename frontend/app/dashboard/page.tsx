'use client'

import { useState } from 'react'
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
  BarChart3
} from 'lucide-react'

export default function DashboardPage() {
  const [isProcessing, setIsProcessing] = useState(false)

  const stats = {
    totalVideos: 24,
    processingVideos: 3,
    publishedVideos: 21,
    totalViews: 156780,
    totalRevenue: 2340.50,
    subscriberGrowth: 12.5
  }

  const recentVideos = [
    {
      id: 1,
      title: "5 AI Tools That Will Replace Your Job (But Make You Rich)",
      status: "published",
      views: 45230,
      likes: 3420,
      comments: 156,
      publishedAt: "2 hours ago",
      thumbnail: "/api/placeholder/320/180"
    },
    {
      id: 2,
      title: "ChatGPT vs Claude: Which AI is Better for Content Creation?",
      status: "processing",
      stage: "video_assembly",
      progress: 75,
      publishedAt: null
    },
    {
      id: 3,
      title: "How I Made $10K with AI Automation (Step-by-Step)",
      status: "published",
      views: 23140,
      likes: 1890,
      comments: 89,
      publishedAt: "1 day ago",
      thumbnail: "/api/placeholder/320/180"
    }
  ]

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
                onClick={() => setIsProcessing(true)}
                disabled={isProcessing}
              >
                <Plus className="w-4 h-4 mr-2" />
                {isProcessing ? 'Creating...' : 'Create Video'}
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
          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total Videos</p>
                  <p className="text-2xl font-bold">{stats.totalVideos}</p>
                </div>
                <Video className="w-8 h-8 text-youtube" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Processing</p>
                  <p className="text-2xl font-bold">{stats.processingVideos}</p>
                </div>
                <Clock className="w-8 h-8 text-ai-purple" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Published</p>
                  <p className="text-2xl font-bold">{stats.publishedVideos}</p>
                </div>
                <TrendingUp className="w-8 h-8 text-automation-green" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Total Views</p>
                  <p className="text-2xl font-bold">{stats.totalViews.toLocaleString()}</p>
                </div>
                <Eye className="w-8 h-8 text-youtube" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Revenue</p>
                  <p className="text-2xl font-bold">${stats.totalRevenue.toFixed(2)}</p>
                </div>
                <DollarSign className="w-8 h-8 text-automation-green" />
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-gray-600">Growth</p>
                  <p className="text-2xl font-bold">+{stats.subscriberGrowth}%</p>
                </div>
                <Users className="w-8 h-8 text-ai-purple" />
              </div>
            </CardContent>
          </Card>
        </div>

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
                  {recentVideos.map((video) => (
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
                              {video.views?.toLocaleString()}
                            </span>
                            <span className="flex items-center">
                              <Heart className="w-3 h-3 mr-1" />
                              {video.likes?.toLocaleString()}
                            </span>
                            <span className="flex items-center">
                              <MessageCircle className="w-3 h-3 mr-1" />
                              {video.comments}
                            </span>
                            <span>{video.publishedAt}</span>
                          </div>
                        ) : (
                          <div className="flex items-center space-x-2">
                            <div className="flex-1 bg-gray-200 rounded-full h-2">
                              <div 
                                className="bg-ai-purple h-2 rounded-full transition-all"
                                style={{ width: `${video.progress}%` }}
                              ></div>
                            </div>
                            <span className="text-xs text-gray-600">{video.progress}%</span>
                          </div>
                        )}
                      </div>
                      
                      <div className="text-right">
                        <div className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          video.status === 'published' 
                            ? 'bg-green-100 text-green-800'
                            : 'bg-blue-100 text-blue-800'
                        }`}>
                          {video.status === 'published' ? 'Published' : 'Processing'}
                        </div>
                      </div>
                    </div>
                  ))}
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