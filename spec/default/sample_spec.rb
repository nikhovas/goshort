require 'spec_helper'

describe package('golang') do
  it { should be_installed }
end

describe service('redis-server') do
  it { should be_enabled }
  it { should be_running }
end

describe service('postgresql') do
  it { should be_enabled }
  it { should be_running }
end

describe port(6379) do
  it { should be_listening }
end

describe port(5432) do
  it { should be_listening }
end
